package snake

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

type Snakeable interface {
	PreRun(ctx context.Context, args []string) (context.Context, error)
	Register(context.Context) (*cobra.Command, error) // need a way for us to be able to have the full command like we do when the kubectl command is importable
}

var (
	ErrMissingBinding   = fmt.Errorf("snake.ErrMissingBinding")
	ErrMissingRun       = fmt.Errorf("snake.ErrMissingRun")
	ErrInvalidRun       = fmt.Errorf("snake.ErrInvalidRun")
	ErrInvalidArguments = fmt.Errorf("snake.ErrInvalidArguments")
)

func NewRootCommand(ctx context.Context, snk Snakeable) (context.Context, error) {

	cmd, err := snk.Register(ctx)
	if err != nil {
		panic(err)
	}

	cmd.SilenceErrors = true

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		zctx := cmd.Context()

		zctx = SetActiveCommand(zctx, RootCommandName)
		defer func() {
			zctx = ClearActiveCommand(zctx)
		}()

		if err := cmd.ParseFlags(args); err != nil {
			return HandleErrorByPrintingToConsole(cmd, err)
		}

		zctx, err := snk.PreRun(zctx, args)
		if err != nil {
			return HandleErrorByPrintingToConsole(cmd, err)
		}

		cmd.SetContext(zctx)

		return nil
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		zctx := cmd.Context()
		zctx = SetActiveCommand(zctx, RootCommandName)
		defer func() {
			zctx = ClearActiveCommand(zctx)
		}()

		zctx, err := snk.PreRun(zctx, args)
		if err != nil {
			return HandleErrorByPrintingToConsole(cmd, err)
		}

		cmd.SetContext(zctx)

		return nil
	}

	ctx = SetRootCommand(ctx, cmd)

	ctx = SetNamedCommand(ctx, RootCommandName, cmd)

	return ctx, nil
}

func Assemble(ctx context.Context) *cobra.Command {
	rootcmd := GetRootCommand(ctx)
	if rootcmd == nil {
		return nil
	}

	named := GetAllNamedCommands(ctx)
	if named == nil {
		return nil
	}

	for name, cmd := range named {
		if name == RootCommandName {
			continue
		}

		rootcmd.AddCommand(cmd)
	}

	rootcmd.SetContext(ctx)

	return rootcmd
}

func MustNewCommand(ctx context.Context, name string, snk Snakeable) context.Context {
	ctx, err := NewCommand(ctx, name, snk)
	if err != nil {
		panic(err)
	}
	return ctx
}

func NewCommand(ctx context.Context, name string, snk Snakeable) (context.Context, error) {

	cmd, err := snk.Register(ctx)
	if err != nil {
		return nil, err
	}

	method := getRunMethod(snk)

	tpe, err := validateRunMethod(snk, method)
	if err != nil {
		return nil, err
	}

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {

		zctx := cmd.Context()

		zctx = SetActiveCommand(zctx, name)
		defer func() {
			zctx = ClearActiveCommand(zctx)
		}()

		if err := cmd.ParseFlags(args); err != nil {
			return HandleErrorByPrintingToConsole(cmd, err)
		}

		zctx, err := snk.PreRun(zctx, args)
		if err != nil {
			return HandleErrorByPrintingToConsole(cmd, err)
		}

		cmd.SetContext(zctx)

		return nil
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		zctx := cmd.Context()

		zctx = SetActiveCommand(zctx, name)
		defer func() {
			zctx = ClearActiveCommand(zctx)
		}()

		// if res := GetBindingResolver(zctx); res != nil {
		// 	resolvers = append(resolvers, res)
		// }

		dctx, err := ResolveBindingsFromProvider(zctx, method, getDynamicBindingResolver(zctx))
		if err != nil {
			return HandleErrorByPrintingToConsole(cmd, err)
		}
		zctx = dctx

		cmd.SetContext(zctx)

		if err := callRunMethod(cmd, method, tpe); err != nil {
			return HandleErrorByPrintingToConsole(cmd, err)
		}
		return nil
	}

	if name != "" {
		cmd.Use = name
	}

	ctx = SetNamedCommand(ctx, name, cmd)

	return ctx, nil
}

func WithRootCommand(ctx context.Context, x func(*cobra.Command) error) error {
	root := GetRootCommand(ctx)
	if root == nil {
		return fmt.Errorf("snake.WithRootCommand: no root command found in context")
	}
	return x(root)
}

func WithNamedCommand(ctx context.Context, name string, x func(*cobra.Command) error) error {
	cmd := GetNamedCommand(ctx, name)
	if cmd == nil {
		return fmt.Errorf("snake.WithNamedCommand: no named command found in context")
	}
	return x(cmd)
}

func WithActiveCommand(ctx context.Context, x func(*cobra.Command) error) error {
	cmd := GetActiveNamedCommand(ctx)
	if cmd == nil {
		return fmt.Errorf("snake.WithActiveCommand: no active command found in context")
	}
	return x(cmd)
}
