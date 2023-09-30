package snake

import (
	"reflect"
	"strings"
	"sync"

	"github.com/go-faster/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Ctx struct {
	bindings  map[string]*reflect.Value
	resolvers map[string]Execable
	cmds      map[string]*cobra.Command

	runlock sync.Mutex
}

var (
	ErrMissingResolver = errors.New("missing resolver")
)

func setBindingWithLock[T any](con *Ctx, val T) func() {
	con.runlock.Lock()
	ptr := reflect.ValueOf(val)
	typ := reflect.TypeOf((*T)(nil)).Elem()
	con.bindings[typ.String()] = &ptr
	return func() {
		delete(con.bindings, typ.String())
		con.runlock.Unlock()
	}
}

func (me *Ctx) getFlags(cmd Execable) (*pflag.FlagSet, error) {
	args := listOfArgs(cmd.Method().Type())
	mapa := make(map[string]bool)
	for _, f := range args {
		if r, ok := me.resolvers[f.String()]; ok {
			mapa[f.String()] = true
			for _, arg := range listOfArgs(r.Method().Type()) {
				// we don't recurse here, because we don't want to process flags more than once if they are used by multiple commands
				if _, ok := me.resolvers[arg.String()]; ok {
					mapa[arg.String()] = true
				} else {
					return nil, errors.Wrapf(ErrMissingResolver, "missing resolver for type %q", arg.String())
				}
			}
		} else {
			return nil, errors.Wrapf(ErrMissingResolver, "missing resolver for type %q", f.String())
		}
	}

	flgs := &pflag.FlagSet{}

	for f := range mapa {
		me.resolvers[f].Flags(flgs)
	}

	return flgs, nil
}

func (me *Ctx) getRunArgs(cmd Execable) ([]reflect.Value, error) {
	args := listOfArgs(cmd.Method().Type())
	rargs := make([]reflect.Value, len(args))
	for i, arg := range args {
		// first check if we have a binding for this type
		if bnd, ok := me.bindings[arg.String()]; ok {
			rargs[i] = *bnd
		} else {
			// if not, check if we have a resolver for this type
			if resl, ok := me.resolvers[arg.String()]; ok {
				// if we do, run it, and store the result as a binding
				// its okay to recurse here, because we are saving the result as a binding
				bnd, err := me.run(resl)
				if err != nil {
					return nil, err
				}
				rargs[i] = bnd
				me.bindings[arg.String()] = &bnd
			} else {
				return nil, errors.Wrapf(ErrMissingResolver, "missing resolver for type %q", arg.String())
			}
		}
	}
	return rargs, nil
}

func (me *Ctx) run(cmd Execable) (reflect.Value, error) {

	args, err := me.getRunArgs(cmd)
	if err != nil {
		return reflect.Value{}, err
	}

	out := cmd.Method().Call(args)

	return cmd.HandleResponse(out)
}

var (
	ErrInvalidMethodSignature = errors.New("invalid method signature")
)

type Flagged interface {
	Flags(*pflag.FlagSet)
}

type MethodType interface {
	ValidateResponse([]reflect.Type) error
	Method() reflect.Value
	HandleResponse([]reflect.Value) (reflect.Value, error)
}

type Execable interface {
	Flags(*pflag.FlagSet)
	ValidateResponse([]reflect.Type) error
	Method() reflect.Value
	HandleResponse([]reflect.Value) (reflect.Value, error)
}

type ExecableCommand struct {
	flags  func(*pflag.FlagSet)
	method reflect.Value
}

func (me *ExecableCommand) Flags(flags *pflag.FlagSet) {
	me.flags(flags)
}

func (me *ExecableCommand) Method() reflect.Value {
	return me.method
}

func (me *ExecableCommand) ValidateResponse(out []reflect.Type) error {

	if len(out) != 1 {
		return errors.Wrapf(ErrInvalidMethodSignature, "invalid return signature, expected 1, got %d", len(out))
	}

	if !out[0].Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return errors.Wrapf(ErrInvalidMethodSignature, "invalid return type %q", out[0].String())
	}

	return nil
}

func (me *ExecableCommand) HandleResponse(out []reflect.Value) (reflect.Value, error) {

	if !out[0].IsNil() {
		return reflect.Zero(reflect.TypeOf(reflect.Value{})), out[1].Interface().(error)
	}

	return end_of_chain, nil
}

type ExecableArgument[I any] struct {
	*ExecableCommand
	result I
}

func (me *ExecableArgument[I]) Method() reflect.Value {
	return me.method
}

func (me *ExecableArgument[I]) HandleResponse(out []reflect.Value) (reflect.Value, error) {

	if !out[1].IsNil() {
		return reflect.Zero(reflect.TypeOf((*I)(nil)).Elem()), out[1].Interface().(error)
	}

	if out[0].Type() != reflect.TypeOf(reflect.TypeOf((*I)(nil)).Elem()) {
		panic("invalid return type")
	}

	me.result = out[0].Interface().(I)

	return out[0], nil
}

func (me *ExecableArgument[I]) ValidateResponse(out []reflect.Type) error {

	if len(out) != 2 {
		return errors.Wrapf(ErrInvalidMethodSignature, "invalid return signature, expected 2, got %d", len(out))
	}

	if !out[0].Implements(reflect.TypeOf((*I)(nil)).Elem()) {
		return errors.Wrapf(ErrInvalidMethodSignature, "invalid return type %q", out[0].String())
	}

	if !out[1].Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return errors.Wrapf(ErrInvalidMethodSignature, "invalid return type %q", out[1].String())
	}

	return nil
}

var end_of_chain = reflect.ValueOf("end_of_chain")
var end_of_chain_ptr = &end_of_chain

func commandName(name string) string {
	return "snake_command_" + name
}

func (me *Ctx) commands() []string {
	var cmds []string
	for k := range me.resolvers {
		if strings.HasPrefix(k, "snake_command_") {
			cmds = append(cmds, strings.TrimPrefix(k, "snake_command_"))
		}
	}
	return cmds
}

func groupCommandName(name string) string {
	return "snake_group_" + name
}

func NewCmdContext(con *Ctx, name string, cbra *cobra.Command, method Flagged) (Execable, error) {

	ec := &ExecableCommand{
		flags:  method.Flags,
		method: getRunMethod(method),
	}

	con.cmds[name] = cbra

	con.resolvers[commandName(name)] = ec

	return ec, nil
}

func NewArgContext[I any](con *Ctx, method Flagged) (Execable, error) {

	ec := &ExecableArgument[I]{
		ExecableCommand: &ExecableCommand{
			flags:  method.Flags,
			method: getRunMethod(method),
		},
	}

	con.resolvers[reflect.TypeOf((*I)(nil)).Elem().String()] = ec

	return ec, nil
}
