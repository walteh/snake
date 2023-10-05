package snake

import (
	"reflect"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Ctx struct {
	bindings  map[string]*reflect.Value
	resolvers map[string]Method
	cmds      map[string]Cobrad

	runlock sync.Mutex
}

var (
	root = Ctx{
		bindings:  make(map[string]*reflect.Value),
		resolvers: make(map[string]Method),
		cmds:      make(map[string]Cobrad),
	}
)

type Flagged interface {
	Flags(*pflag.FlagSet)
}

type Cobrad interface {
	Cobra() *cobra.Command
}

func NewArgument[I any](method Flagged) {
	_ = NewArgContext[I](&root, method)
}

func NewCmd[I Cobrad](cmd I, method Flagged) {
	_ = NewCmdContext(&root, cmd, method)
}

func NewCmdContext[I Cobrad](con *Ctx, cbra I, m Flagged) Method {

	ec := &method{
		flags:              m.Flags,
		validationStrategy: commandResponseValidationStrategy,
		responseStrategy:   commandResponseHandleStrategy,
		// name:               prefix_command + name,
		name:   reflect.TypeOf((*I)(nil)).Elem().String(),
		method: getRunMethod(m),
		cmd:    cbra,
	}

	con.runlock.Lock()
	defer con.runlock.Unlock()

	con.resolvers[ec.name] = ec

	return ec
}

func methodName(typ reflect.Type) string {
	return prefix_argument + typ.String()
}

func NewArgContext[I any](con *Ctx, m Flagged) Method {

	ec := &method{
		flags:              m.Flags,
		validationStrategy: validateArgumentResponse[I],
		responseStrategy:   handleArgumentResponse[I],
		name:               methodName(reflect.TypeOf((*I)(nil)).Elem()),
		method:             getRunMethod(m),
	}

	con.runlock.Lock()
	defer con.runlock.Unlock()

	con.resolvers[ec.name] = ec

	return ec
}
