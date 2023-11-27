package scobra

import (
	"os"
	"strings"

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
}

func (me *CS) Decorate(self SCobra, snk sbind.Snake, inputs []sbind.Input) error {

	cmd := self.Command()

	name := cmd.Name()

	oldRunE := cmd.RunE

	for _, v := range inputs {
		flgs := cmd.Flags()

		// a := cmd.Flags().Lookup(strings.ToLower(v.Name()))
		// if a == nil {
		// 	a = cmd.PersistentFlags().Lookup(strings.ToLower(v.Name()))
		// }

		if v.Shared() {
			flgs = me.PersistentFlags()
		} else {
			if cmd.Flags().Lookup(strings.ToLower(v.Name())) != nil {
				// if this is the same object, then the user is trying to override the flag, so we let them
				continue
			}
		}

		// pp.Println("looking up flag", v.Name(), v.Usage(), v.M(), self, a)

		// if a != nil {

		// this was working and would be cool, but it's not predictable enough
		// 	if v.M() == self {
		// 		// if this is the same object, then the user is trying to override the flag, so we let them
		// 		continue
		// 	}

		// 	return errors.Errorf("multiple flags named %q resolved for %q", a.Name, reflect.ValueOf(v.M()).Elem().Type().String())
		// }

		// pp.Println("adding flag", v.Name(), v.Usage(), v.M(), self)

		switch t := v.(type) {
		case *sbind.StringInput:
			flgs.StringVar(t.Value(), strings.ToLower(v.Name()), t.Default(), v.Usage())
		case *sbind.BoolInput:
			flgs.BoolVar(t.Value(), strings.ToLower(v.Name()), t.Default(), v.Usage())
		case *sbind.IntInput:
			flgs.IntVar(t.Value(), strings.ToLower(v.Name()), t.Default(), v.Usage())
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
