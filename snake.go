package snake

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

type Snakeable interface {
	ParseArguments(ctx context.Context, cmd *cobra.Command, args []string) error
	BuildCommand(ctx context.Context) *cobra.Command
}

var (
	ErrMissingBinding   = fmt.Errorf("snake.ErrMissingBinding")
	ErrMissingRun       = fmt.Errorf("snake.ErrMissingRun")
	ErrInvalidRun       = fmt.Errorf("snake.ErrInvalidRun")
	ErrInvalidArguments = fmt.Errorf("snake.ErrInvalidArguments")
)

func NewRootCommand(ctx context.Context, snk Snakeable) context.Context {

	cmd := snk.BuildCommand(ctx)

	cmd.SilenceErrors = true

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := cmd.ParseFlags(args); err != nil {
			return HandleErrorByPrintingToConsole(cmd, err)
		}

		err := snk.ParseArguments(cmd.Context(), cmd, args)
		if err != nil {
			return HandleErrorByPrintingToConsole(cmd, err)
		}

		return nil
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := snk.ParseArguments(cmd.Context(), cmd, args)
		if err != nil {
			return HandleErrorByPrintingToConsole(cmd, err)
		}
		return nil
	}

	// if prov, ok := snk.(BindingResolver); ok {
	// 	ctx = SetBindingResolver(ctx, prov)
	// }

	ctx = SetRootCommand(ctx, cmd)

	return ctx

}

func MustNewCommand(ctx context.Context, name string, snk Snakeable) context.Context {
	ctx, err := NewCommand(ctx, name, snk)
	if err != nil {
		panic(err)
	}
	return ctx
}

func NewCommand(ctx context.Context, name string, snk Snakeable) (context.Context, error) {

	rootcmd := GetRootCommand(ctx)
	if rootcmd == nil {
		return nil, fmt.Errorf("snake.NewCommand: cannot create a new command when a root command has already been created")
	}

	cmd := snk.BuildCommand(ctx)

	method := getRunMethod(snk)

	tpe, err := validateRunMethod(snk, method)
	if err != nil {
		return nil, err
	}

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {

		if err := cmd.ParseFlags(args); err != nil {
			return HandleErrorByPrintingToConsole(cmd, err)
		}

		err := snk.ParseArguments(cmd.Context(), cmd, args)
		if err != nil {
			return HandleErrorByPrintingToConsole(cmd, err)
		}

		return nil
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		zctx := cmd.Context()

		resolvers := []BindingResolver{}

		if res := GetBindingResolver(zctx); res != nil {
			resolvers = append(resolvers, res)
		}

		if res := getDynamicBindingResolver(zctx); res != nil {
			resolvers = append(resolvers, res)
		}

		if len(resolvers) > 0 {
			dctx, err := ResolveBindingsFromProvider(zctx, method, resolvers...)
			if err != nil {
				return HandleErrorByPrintingToConsole(cmd, err)
			}
			zctx = dctx
		}

		cmd.SetContext(zctx)

		if err := callRunMethod(cmd, method, tpe); err != nil {
			return HandleErrorByPrintingToConsole(cmd, err)
		}
		return nil
	}

	if name != "" {
		cmd.Use = name
	}

	rootcmd.AddCommand(cmd)

	return ctx, nil
}

func UseRootCommand(ctx context.Context, x func(*cobra.Command) error) error {
	root := GetRootCommand(ctx)
	if root == nil {
		return fmt.Errorf("snake.UseRootCommand: no root command found in context")
	}
	return x(root)
}
