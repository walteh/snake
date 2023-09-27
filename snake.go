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

func NewRootCommand(ctx context.Context, snk Snakeable) *cobra.Command {

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

	if prov, ok := snk.(BindingResolver); ok {
		ctx = SetBindingResolver(ctx, prov)
	}

	cmd.SetContext(ctx)

	return cmd

}

func MustNewCommand(cbra *cobra.Command, name string, snk Snakeable) {
	err := NewCommand(cbra, name, snk)
	if err != nil {
		panic(err)
	}
}

func NewCommand(cbra *cobra.Command, name string, snk Snakeable) error {

	ctx := cbra.Context()

	cmd := snk.BuildCommand(ctx)

	method := getRunMethod(snk)

	tpe, err := validateRunMethod(snk, method)
	if err != nil {
		return err
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

		resolvers := []BindingResolver{}

		if res := GetBindingResolver(ctx); res != nil {
			resolvers = append(resolvers, res)
		}

		if res := getDynamicBindingResolver(ctx); res != nil {
			resolvers = append(resolvers, res)
		}

		if len(resolvers) > 0 {
			if err := ResolveBindingsFromProvider(cmd, method, resolvers...); err != nil {
				return HandleErrorByPrintingToConsole(cmd, err)
			}
		}

		if err := callRunMethod(cmd, method, tpe); err != nil {
			return HandleErrorByPrintingToConsole(cmd, err)
		}
		return nil
	}

	if name != "" {
		cmd.Use = name
	}

	cbra.AddCommand(cmd)

	return nil
}
