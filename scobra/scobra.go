package scobra

import (
	"context"
	"os"
	"strings"

	"github.com/go-faster/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/walteh/snake"
)

type CS struct {
	*cobra.Command
}

type SCobra interface {
	Command() *cobra.Command
	Name() string
	// sbind.NamedMethod
}

func NewCommandResolver(s SCobra) snake.TypedResolver[SCobra] {
	return snake.MustGetRunMethod(s)
}

func (me *CS) ManagedResolvers(_ context.Context) []snake.Resolver {
	return []snake.Resolver{
		snake.NewNoopMethod[*cobra.Command](),
		snake.NewNoopMethod[[]string](),
	}
}

func (me *CS) Decorate(ctx context.Context, self SCobra, snk snake.Snake, inputs []snake.Input) error {

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
		case *snake.StringEnumInput:
			flgs.Var(NewWrappedEnum(t), v.Name(), t.Usage())
		case *snake.StringInput:
			flgs.StringVar(t.Value(), v.Name(), t.Default(), t.Usage())
		case *snake.BoolInput:
			flgs.BoolVar(t.Value(), v.Name(), t.Default(), t.Usage())
		case *snake.IntInput:
			flgs.IntVar(t.Value(), v.Name(), t.Default(), t.Usage())
		case *snake.StringArrayInput:
			flgs.StringSliceVar(t.Value(), v.Name(), t.Default(), t.Usage())
		case *snake.IntArrayInput:
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
		binder := snake.NewBinder()

		snake.SetBinding(binder, cmd)
		snake.SetBinding(binder, args)

		outhand := NewOutputHandler(cmd)

		err := snake.RunResolvingArguments(outhand, snk.Resolve, name, binder)
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

func (me *CS) OnSnakeInit(ctx context.Context, snk snake.Snake) error {

	me.RunE = func(cmd *cobra.Command, args []string) error {
		binder := snake.NewBinder()

		snake.SetBinding(binder, cmd)
		snake.SetBinding(binder, args)

		outhand := NewOutputHandler(cmd)

		err := snake.RunResolvingArguments(outhand, snk.Resolve, "root", binder)
		if err != nil {
			return HandleErrorByPrintingToConsole(cmd, err)
		}
		return nil
	}

	return nil
}

func NewCobraSnake(root *cobra.Command) (*CS, error) {

	if root == nil {
		root = &cobra.Command{}
	}

	me := &CS{root}

	root.SilenceUsage = true

	return me, nil
}
