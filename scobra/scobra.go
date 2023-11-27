package scobra

import (
	"os"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/walteh/snake/sbind"
)

type CS struct {
	*cobra.Command
}

func (me *CS) Decorate(self SCobra, snk sbind.Snake) error {

	cmd := self.Command()

	name := cmd.Name()

	oldRunE := cmd.RunE

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

type SCobra interface {
	Command() *cobra.Command
}

type Flagged interface {
	Flags(*pflag.FlagSet)
}

type NewSCobraOpts struct {
	Commands  []SCobra
	Resolvers []Flagged
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

	root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		name := cmd.Name()

		if flgs, err := sbind.FlagsFor(name, snk.Resolve); err != nil {
			return err
		} else {
			fs := pflag.NewFlagSet(name, pflag.ContinueOnError)

			procd := make(map[any]bool, 0)

			for _, f := range flgs {
				if z := snk.Resolve(f); reflect.ValueOf(z).IsNil() {
					return errors.Errorf("missing resolver for %q", f)
				} else {
					if z, ok := z.(Flagged); ok && !procd[z] {
						procd[z] = true
						z.Flags(fs)
					}
				}
			}

			fs.VisitAll(func(f *pflag.Flag) {
				// if globalFlags != nil && globalFlags.Lookup(f.Name) != nil {
				// 	return
				// }
				cmd.Flags().AddFlag(f)
			})
		}

		return nil
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
