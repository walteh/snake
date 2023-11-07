package snake

import (
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Snake struct {
	bindings  map[string]*reflect.Value
	resolvers map[string]Method
	root      *cobra.Command
	runlock   sync.Mutex
}

type Flagged interface {
	Flags(*pflag.FlagSet)
}

type Cobrad interface {
	Cobra() *cobra.Command
}

type NewSnakeOpts struct {
	Root                       *cobra.Command
	Commands                   []Method
	Resolvers                  []Method
	GlobalContextResolverFlags bool
}

func attachMethod(me *Snake, exer Method) (*cobra.Command, error) {

	if exer.Command() == nil {
		return nil, nil
	}

	cmd := exer.Command().Cobra()

	if flgs, err := FlagsFor(exer.Name(), func(s string) Method {
		return me.resolvers[s]
	}); err != nil {
		return nil, err
	} else {
		cmd.Flags().AddFlagSet(flgs)
	}

	err := exer.ValidateResponse()
	if err != nil {
		return nil, err
	}

	oldRunE := cmd.RunE

	// if a flag is not set, we check the environment for "cmd_name_arg_name"

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if f.Changed {
				return
			}
			val := strings.ToUpper(me.root.Name() + "_" + strings.ReplaceAll(f.Name, "-", "_"))
			envvar := os.Getenv(val)
			if envvar == "" {
				return
			}
			err := f.Value.Set(envvar)
			if err != nil {
				return
			}
		})
		return nil
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		defer setBindingWithLock(me, cmd)()
		defer setBindingWithLock(me, args)()

		err := runResolvingArguments(exer.Name(), func(s string) IsRunnable {
			return me.resolvers[s]
		}, me.bindings)
		if err != nil {
			return HandleErrorByPrintingToConsole(cmd, err)
		}
		if oldRunE != nil {
			err := oldRunE(cmd, args)
			if err != nil {
				return HandleErrorByPrintingToConsole(cmd, err)
			}
		}
		return nil
	}

	return cmd, nil

}

func NewSnake(opts *NewSnakeOpts) (*cobra.Command, error) {

	root := opts.Root

	if root == nil {
		root = &cobra.Command{}
	}

	snk := &Snake{
		bindings:  make(map[string]*reflect.Value),
		resolvers: make(map[string]Method),
		root:      root,
	}

	for _, v := range opts.Resolvers {
		snk.resolvers[v.Name()] = v

		if opts.GlobalContextResolverFlags && v.IsContextResolver() {
			v.Flags(root.PersistentFlags())
		}
	}

	for _, v := range opts.Commands {
		snk.resolvers[v.Name()] = v
	}

	// these will always be overwritten in the RunE function
	snk.resolvers["*cobra.Command"] = NewArgumentMethod[*cobra.Command](&inlineResolver[*cobra.Command]{
		flagFunc: func(*pflag.FlagSet) {},
		runFunc: func() (*cobra.Command, error) {
			return &cobra.Command{}, nil
		},
	})

	snk.resolvers["[]string"] = NewArgumentMethod[[]string](&inlineResolver[[]string]{
		flagFunc: func(*pflag.FlagSet) {},
		runFunc: func() ([]string, error) {
			return []string{}, nil
		},
	})

	for _, exer := range snk.resolvers {
		if cmd, err := attachMethod(snk, exer); err != nil {
			return nil, err
		} else if cmd != nil {
			root.AddCommand(cmd)
		}
	}

	return root, nil
}

func NewCommandMethod[I Cobrad](cbra I) Method {

	ec := &method{
		flags:              func(*pflag.FlagSet) {},
		validationStrategy: commandResponseValidationStrategy,
		responseStrategy:   commandResponseHandleStrategy,
		name:               reflect.TypeOf((*I)(nil)).Elem().String(),
		method:             getRunMethod(cbra),
		cmd:                cbra,
	}

	if flg, ok := any(cbra).(Flagged); ok {
		ec.flags = flg.Flags
	}

	return ec
}

func NewArgumentMethod[I any](m Flagged) Method {

	ec := &method{
		flags:              m.Flags,
		validationStrategy: validateArgumentResponse[I],
		responseStrategy:   handleArgumentResponse[I],
		name:               reflect.TypeOf((*I)(nil)).Elem().String(),
		method:             getRunMethod(m),
	}

	return ec
}
