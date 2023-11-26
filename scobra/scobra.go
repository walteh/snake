package scobra

import (
	"os"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/walteh/snake"
	"github.com/walteh/snake/sbind"
)

func attachMethod(me *snake.Snake, cmd *cobra.Command, name string, globalFlags *pflag.FlagSet) (*cobra.Command, error) {

	if cmd == nil {
		return nil, nil
	}

	if flgs, err := sbind.FlagsFor[pflag.FlagSet, Method](name, func(s string) Method {
		return me.resolvers[s]
	}); err != nil {
		return nil, err
	} else {
		flgs.VisitAll(func(f *pflag.Flag) {
			if globalFlags != nil && globalFlags.Lookup(f.Name) != nil {
				return
			}
			cmd.Flags().AddFlag(f)
		})
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

		err := sbind.RunResolvingArguments(name, func(s string) sbind.IsRunnable {
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
		for _, n := range v.Names() {
			snk.resolvers[n] = v
		}

		if opts.GlobalContextResolverFlags && v.IsContextResolver() {
			v.Flags(root.PersistentFlags())
		}
	}

	for _, v := range opts.Commands {
		for _, n := range v.Names() {
			snk.resolvers[n] = v
		}
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
