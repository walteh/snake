package snake

import (
	"context"
	"reflect"
	"sync"

	"github.com/go-faster/errors"
	"github.com/spf13/cobra"
	"github.com/walteh/snake/sbind"
)

type Snake struct {
	bindings  map[string]*reflect.Value
	resolvers map[string]Method
	root      *cobra.Command
	runlock   sync.Mutex
}

type Cobrad interface {
	Cobra() *cobra.Command
}

type NewSnakeOpts[F any] struct {
	Root                       *cobra.Command
	Commands                   []Method
	Resolvers                  []Method
	GlobalContextResolverFlags bool
}

func NewCommandMethod[I Cobrad](cbra I) Method {

	ec := &method{
		// flags:              func(*pflag.FlagSet) {},
		validationStrategy: commandResponseValidationStrategy,
		responseStrategy:   commandResponseHandleStrategy,
		names:              []string{reflect.TypeOf((*I)(nil)).Elem().String()},
		method:             sbind.GetRunMethod(cbra),
		// cmd:                cbra,
	}

	if flg, ok := any(cbra).(Flagged); ok {
		ec.setFlag = flg.SetFlag
		ec.getFlag = flg.GetFlag
	}

	return ec
}

func NewArgumentMethod[A any](m Flagged) Method {

	ec := &method{
		setFlag:            m.SetFlag,
		getFlag:            m.GetFlag,
		validationStrategy: validate1ArgumentResponse[A],
		responseStrategy:   handle1ArgumentResponse[A],
		names:              namesBuilder((*A)(nil)),
		method:             sbind.GetRunMethod(m),
	}

	return ec
}

func New2ArgumentMethod[A any, B any](m Flagged) Method {

	ec := &method{
		setFlag:            m.SetFlag,
		getFlag:            m.GetFlag,
		validationStrategy: validate2ArgumentResponse[A, B],
		responseStrategy:   handle2ArgumentResponse[A, B],
		names:              namesBuilder((*A)(nil), (*B)(nil)),
		method:             sbind.GetRunMethod(m),
	}

	return ec
}

func New3ArgumentMethod[A any, B any, C any](m Flagged) Method {

	ec := &method{
		setFlag:            m.SetFlag,
		getFlag:            m.GetFlag,
		validationStrategy: validate3ArgumentResponse[A, B, C],
		responseStrategy:   handle3ArgumentResponse[A, B, C],
		names:              namesBuilder((*A)(nil), (*B)(nil), (*C)(nil)),
		method:             sbind.GetRunMethod(m),
	}

	return ec
}

func namesBuilder(inter ...any) []string {

	var names []string
	for _, v := range inter {
		names = append(names, reflect.TypeOf(v).Elem().String())
	}
	return names
}

type fakeCobra struct {
}

func (me *fakeCobra) Cobra() *cobra.Command {
	return &cobra.Command{}
}

func (me *fakeCobra) Run(cmd *cobra.Command) error {
	return errors.Errorf("no method found for %q", cmd.Name())
}

type fakeCobraWithContext struct {
	internal fakeCobra
}

func (me *fakeCobraWithContext) Cobra() *cobra.Command {
	return me.internal.Cobra()
}

func (me *fakeCobraWithContext) Run(_ context.Context, cmd *cobra.Command) error {
	return me.internal.Run(cmd)
}
