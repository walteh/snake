package scobra

import (
	"context"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/go-faster/errors"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/walteh/snake"
)

var (
	_ snake.SnakeImplementation[SCobra] = &CobraSnake{}
)

type CobraSnake struct {
	RootCommand *cobra.Command
}

type SCobra interface {
	Command() *cobra.Command
	Name() string
}

func NewCommandResolver(s SCobra) snake.TypedResolver[SCobra] {
	return snake.MustGetTypedResolver(s)
}

func (me *CobraSnake) ManagedResolvers(_ context.Context) []snake.Resolver {
	return []snake.Resolver{
		snake.NewNoopMethod[*cobra.Command](),
		snake.NewNoopMethod[[]string](),
	}
}

func applyInputToFlags(input snake.Input, flgs *pflag.FlagSet) error {
	switch t := input.(type) {
	case *snake.StringEnumInput:
		flgs.Var(NewWrappedEnum(t), input.Name(), t.Usage())
	case *snake.StringInput:
		flgs.StringVar(t.Value(), input.Name(), t.Default(), t.Usage())
	case *snake.BoolInput:
		flgs.BoolVar(t.Value(), input.Name(), t.Default(), t.Usage())
	case *snake.IntInput:
		flgs.IntVar(t.Value(), input.Name(), t.Default(), t.Usage())
	case *snake.StringArrayInput:
		flgs.StringSliceVar(t.Value(), input.Name(), t.Default(), t.Usage())
	case *snake.IntArrayInput:
		flgs.IntSliceVar(t.Value(), input.Name(), t.Default(), t.Usage())
	case *snake.DurationInput:
		flgs.DurationVar(t.Value(), input.Name(), t.Default(), t.Usage())
	default:
		return errors.Errorf("unknown input type %T", t)
	}
	return nil
}

func (me *CobraSnake) Decorate(ctx context.Context, self SCobra, snk snake.Snake, inputs []snake.Input, mw []snake.Middleware) error {

	cmd := self.Command()

	name := cmd.Name()

	oldRunE := cmd.RunE

	for _, v := range inputs {
		flgs := cmd.Flags()

		if v.Shared() {
			flgs = me.RootCommand.PersistentFlags()
		} else {
			if cmd.Flags().Lookup(v.Name()) != nil {
				// if this is the same object, then the user is trying to override the flag, so we let them
				continue
			}
		}

		err := applyInputToFlags(v, flgs)
		if err != nil {
			return err
		}
	}

	// if a flag is not set, we check the environment for "cmd_name_arg_name"
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if f.Changed {
				return
			}
			val := strings.ToUpper(me.RootCommand.Name() + "_" + strings.ReplaceAll(f.Name, "-", "_"))
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

		err := snake.RunResolvingArguments(outhand, snk.Resolve, name, binder, mw...)
		if err != nil {
			return err
		}
		if oldRunE != nil {
			err := oldRunE(cmd, args)
			if err != nil {
				return err
			}
		}
		return nil
	}

	me.RootCommand.AddCommand(cmd)

	return nil
}

func (me *CobraSnake) OnSnakeInit(ctx context.Context, snk snake.Snake) error {

	me.RootCommand.RunE = func(cmd *cobra.Command, args []string) error {
		binder := snake.NewBinder()

		snake.SetBinding(binder, cmd)
		snake.SetBinding(binder, args)

		outhand := NewOutputHandler(cmd)

		err := snake.RunResolvingArguments(outhand, snk.Resolve, "root", binder)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

var _ snake.EnumResolverFunc = (*CobraSnake)(nil).ResolveEnum

func (me *CobraSnake) ResolveEnum(typ string, opts []string) (string, error) {
	prompt := promptui.Select{
		Label: "Select " + typ,
		Items: opts,
	}

	_, result, err := prompt.Run()

	if err != nil {
		return "", err
	}

	if result == "" {
		return "", errors.Errorf("invalid %q", typ)
	}

	return result, nil

}

func (me *CobraSnake) ProvideContextResolver() snake.Resolver {
	return snake.MustGetResolverFor[context.Context](&ContextResolver{})
}

func NewCobraSnake(ctx context.Context, root *cobra.Command) *CobraSnake {

	if root == nil {
		root = &cobra.Command{}
	}

	me := &CobraSnake{root}

	str, err := DecorateTemplate(ctx, root, &DecorateOptions{
		Headings: color.New(color.FgCyan, color.Bold),
		ExecName: color.New(color.FgHiGreen, color.Bold),
		Commands: color.New(color.FgHiRed, color.Faint),
	})
	if err != nil {
		panic(err)
	}

	root.SetUsageTemplate(str)

	root.SilenceUsage = true

	return me
}

func ExecuteHandlingError(ctx context.Context, cmd *CobraSnake) {
	err := HandleErrorByPrintingToConsole(cmd.RootCommand, cmd.RootCommand.ExecuteContext(ctx))
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
