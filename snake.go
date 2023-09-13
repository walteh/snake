package snake

import (
	"context"
	"fmt"
	"regexp"

	"github.com/fatih/color"
	"github.com/go-faster/errors"
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

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := cmd.ParseFlags(args); err != nil {
			return wrapPrintHandleError(cmd, err)
		}

		err := snk.ParseArguments(cmd.Context(), cmd, args)
		if err != nil {
			return wrapPrintHandleError(cmd, err)
		}

		return nil
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := snk.ParseArguments(cmd.Context(), cmd, args)
		if err != nil {
			return wrapPrintHandleError(cmd, err)
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
			return wrapPrintHandleError(cmd, err)
		}

		err := snk.ParseArguments(cmd.Context(), cmd, args)
		if err != nil {
			return wrapPrintHandleError(cmd, err)
		}
		return nil
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := callRunMethod(cmd, method, tpe)
		if err != nil {
			return wrapPrintHandleError(cmd, err)
		}
		return nil
	}

	if name != "" {
		cmd.Use = name
	}

	cbra.AddCommand(cmd)

	return nil
}

func wrapPrintHandleError(cmd *cobra.Command, err error) error {
	cmd.Println(FormatError(cmd, err))
	return errors.Wrap(ErrHandled, err.Error())
}

func FormatError(cmd *cobra.Command, err error) string {

	n := color.New(color.FgHiRed).Sprint(cmd.Name())
	cmd.VisitParents(func(cmd *cobra.Command) {
		if cmd.Name() != "" {
			n = cmd.Name() + " " + n
		}
	})
	caller := ""
	if frm, ok := errors.Cause(err); ok {
		_, filestr, linestr := frm.Location()
		caller = FormatCaller(filestr, linestr)
		caller = caller + " - "

	}
	str := fmt.Sprintf("%+s", err)
	prev := ""
	// replace any string that contains "*.Err" with a bold red version using regex
	str = regexp.MustCompile(`\S+\.Err\S*`).ReplaceAllStringFunc(str, func(s string) string {
		prev += color.New(color.FgRed, color.Bold).Sprint(s) + " -> "
		return ""
	})

	return fmt.Sprintf("%s - %s - %s%s%s\n", color.New(color.FgRed, color.Bold).Sprint("ERROR"), n, caller, prev, color.New(color.FgRed).Sprint(str))
}
