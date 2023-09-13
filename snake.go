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
	ErrHandled          = fmt.Errorf("snake.ErrHandled")
)

func NewRootCommand(ctx context.Context, snk Snakeable) *cobra.Command {

	cmd := snk.BuildCommand(ctx)

	// adapted from https://github.com/spf13/cobra/issues/914#issuecomment-548411337
	cmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "[%s] (error) %+v\n", cmd.Name(), err)
		return ErrHandled
	})

	cmd.SilenceErrors = true

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := cmd.ParseFlags(args); err != nil {
			return err
		}

		err := snk.ParseArguments(cmd.Context(), cmd, args)
		if err != nil {
			return err
		}

		return nil
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := snk.ParseArguments(cmd.Context(), cmd, args)
		if err != nil {
			return err
		}
		return nil
	}

	cmd.SetContext(ctx)

	return cmd

}

func NewGroup(ctx context.Context, cmd *cobra.Command, name string, description string) *cobra.Command {

	grp := &cobra.Command{
		Use:   name,
		Short: description,
	}

	cmd.AddCommand(grp)

	return grp
}

func MustNewCommand(ctx context.Context, cbra *cobra.Command, name string, snk Snakeable) {
	err := NewCommand(ctx, cbra, name, snk)
	if err != nil {
		panic(err)
	}
}

func NewCommand(ctx context.Context, cbra *cobra.Command, name string, snk Snakeable) error {

	cmd := snk.BuildCommand(ctx)

	method := getRunMethod(snk)

	tpe, err := validateRunMethod(snk, method)
	if err != nil {
		return err
	}

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {

		if err := cmd.ParseFlags(args); err != nil {
			return err
		}
		err := snk.ParseArguments(cmd.Context(), cmd, args)
		if err != nil {
			return err
		}
		return nil
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return callRunMethod(cmd, method, tpe)
	}

	if name != "" {
		cmd.Use = name
	}

	cbra.AddCommand(cmd)

	return nil
}
