package scobra

import (
	"os"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/walteh/snake"
	"github.com/walteh/snake/sbind"
)

type CobraMethod interface {
	Command() *cobra.Command
	Flags(*pflag.FlagSet)
}

type cobraSnake struct {
	bindings *sbind.Binder
	internal sbind.Snake
}

func (c *cobraSnake) attachMethod(me snake.Snake, cmd *cobra.Command, name string, globalFlags *pflag.FlagSet) (*cobra.Command, error) {

	if cmd == nil {
		return nil, nil
	}

	if flgs, err := sbind.FlagsFor(name, me.Resolve); err != nil {
		return nil, err
	} else {
		fs := pflag.NewFlagSet(name, pflag.ContinueOnError)

		procd := make(map[any]bool, 0)

		for _, f := range flgs {
			if z := me.Resolve(f); reflect.ValueOf(z).IsNil() {
				return nil, errors.Errorf("missing resolver for %q", f)
			} else {
				if z, ok := any(z).(CobraMethod); ok && !procd[z] {
					procd[z] = true
					z.Flags(fs)
				}
			}
		}

		fs.VisitAll(func(f *pflag.Flag) {
			if globalFlags != nil && globalFlags.Lookup(f.Name) != nil {
				return
			}
			cmd.Flags().AddFlag(f)
		})
	}

	return cmd, nil

}

func (c *cobraSnake) Decorate(cmd *cobra.Command) error {

	name := cmd.Name()

	oldRunE := cmd.RunE

	// if a flag is not set, we check the environment for "cmd_name_arg_name"

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if f.Changed {
				return
			}
			val := strings.ToUpper(cmd.CommandPath() + "_" + strings.ReplaceAll(f.Name, "-", "_"))
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
		defer sbind.SetBindingWithLock(c.bindings, cmd)()
		defer sbind.SetBindingWithLock(c.bindings, args)()

		err := sbind.RunResolvingArguments(name, c.internal.Resolve)
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

	return nil
}

func NewCobraSnake(root *cobra.Command, opts *sbind.NewSnakeOpts) (*cobra.Command, error) {

	if root == nil {
		root = &cobra.Command{}
	}

	snk, err := sbind.NewSnake(opts)
	if err != nil {
		return nil, err
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
		if exer.Command() == nil {
			continue
		}
		name := exer.Names()[0]
		if cmd, err := attachMethod(snk, exer.Command().Cobra(), name, root.PersistentFlags()); err != nil {
			return nil, err
		} else if cmd != nil {
			err := exer.ValidateResponse()
			if err != nil {
				return nil, err
			}
			root.AddCommand(cmd)
		}
	}

	if opts.GlobalContextResolverFlags {
		// this will force the context to be resolved before any command is run
		snk.resolvers["root"] = NewCommandMethod(&fakeCobraWithContext{})
	} else {
		snk.resolvers["root"] = NewCommandMethod(&fakeCobra{})
	}

	root.RunE = func(cmd *cobra.Command, args []string) error {
		err := sbind.RunResolvingArguments("root", func(s string) sbind.IsRunnable {
			return snk.resolvers[s]
		}, snk.bindings)
		if err != nil {
			return HandleErrorByPrintingToConsole(cmd, err)
		}
		return nil
	}

	root.SilenceUsage = true

	return root, nil
}
