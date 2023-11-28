package scobra

import (
	"os"
	"strings"

	"github.com/go-faster/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/walteh/snake/sbind"
)

type CS struct {
	*cobra.Command
}

type SCobra interface {
	Command() *cobra.Command
}

type NewSCobraOpts struct {
	Commands  []SCobra
	Resolvers []sbind.Method
	Enums     []sbind.EnumOption
}

func (me *CS) Decorate(self SCobra, snk sbind.Snake, inputs []sbind.Input) error {

	cmd := self.Command()

	name := cmd.Name()

	oldRunE := cmd.RunE

	for _, v := range inputs {
		flgs := cmd.Flags()

		if v.Shared() {
			flgs = me.PersistentFlags()
		} else {
			if cmd.Flags().Lookup(v.Name()) != nil {
				// if this is the same object, then the user is trying to override the flag, so we let them
				continue
			}
		}

		switch t := v.(type) {
		case *sbind.StringEnumInput:
			d, err := NewWrappedEnum(t.Default(), t.Value(), t.Options()...)
			if err != nil {
				return err
			}
			flgs.Var(d, v.Name(), t.Usage())
		case *sbind.IntEnumInput:
			d, err := NewWrappedEnum(t.Default(), t.Value(), t.Options()...)
			if err != nil {
				return err
			}
			flgs.Var(d, v.Name(), t.Usage())
		case *sbind.StringInput:
			flgs.StringVar(t.Value(), v.Name(), t.Default(), t.Usage())
		case *sbind.BoolInput:
			flgs.BoolVar(t.Value(), v.Name(), t.Default(), t.Usage())
		case *sbind.IntInput:
			flgs.IntVar(t.Value(), v.Name(), t.Default(), t.Usage())
		case *sbind.StringArrayInput:
			flgs.StringSliceVar(t.Value(), v.Name(), t.Default(), t.Usage())
		case *sbind.IntArrayInput:
			flgs.IntSliceVar(t.Value(), v.Name(), t.Default(), t.Usage())
		default:
			return errors.Errorf("unknown input type %T", t)
		}
	}

	// if a flag is not set, we check the environment for "cmd_name_arg_name"
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if f.Changed {
				return
			}
			val := strings.ToUpper(me.Name() + "_" + strings.ReplaceAll(f.Name, "-", "_"))
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
		defer sbind.SetBindingWithLock(snk.Binder(), cmd)()
		defer sbind.SetBindingWithLock(snk.Binder(), args)()

		err := sbind.RunResolvingArguments(name, snk.Resolve)
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

	me.AddCommand(cmd)

	return nil
}

func NewCobraSnake(root *cobra.Command, opts *NewSCobraOpts) (*cobra.Command, error) {

	if root == nil {
		root = &cobra.Command{}
	}

	me := &CS{root}

	opts2 := &sbind.NewSnakeOpts{
		Resolvers:      make([]sbind.Method, 0),
		NamedResolvers: map[string]sbind.Method{},
		Enums:          opts.Enums,
	}

	for _, v := range opts.Commands {
		opts2.NamedResolvers[v.Command().Name()] = v
	}

	for _, v := range opts.Resolvers {
		opts2.Resolvers = append(opts2.Resolvers, v)
	}

	// these will always be overwritten in the RunE function
	opts2.Resolvers = append(opts2.Resolvers, sbind.NewNoopMethod[*cobra.Command]())
	opts2.Resolvers = append(opts2.Resolvers, sbind.NewNoopMethod[[]string]())

	snk, err := sbind.NewSnake(opts2, me)
	if err != nil {
		return nil, err
	}

	root.RunE = func(cmd *cobra.Command, args []string) error {
		err := sbind.RunResolvingArguments("root", snk.Resolve)
		if err != nil {
			return HandleErrorByPrintingToConsole(cmd, err)
		}
		return nil
	}

	root.SilenceUsage = true

	return root, nil
}
