package snake

import (
	"context"
	"fmt"
	"reflect"

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

	nc := &NamedCommand{
		cmd:        cmd,
		method:     reflect.ValueOf(func() {}),
		methodType: reflect.TypeOf(func() {}),
		ptr:        nil,
	}

	ctx = SetRootCommand(ctx, nc)

	ctx = SetNamedCommand(ctx, RootCommandName, nc)

	return ctx, nil
}

type NamedCommand struct {
	cmd        *cobra.Command
	method     reflect.Value
	methodType reflect.Type
	ptr        any
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

	rootcmd := GetRootCommand(ctx)
	if rootcmd == nil {
		return nil, fmt.Errorf("snake.NewCommand: no root command found in context")
	}

	getRootCommandCtx := func() context.Context {
		if c, ok := rootcmd.ptr.(context.Context); ok {
			return c
		}
		return ctx
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

		dctx, err := ResolveBindingsFromProvider(getRootCommandCtx(), method)
		if err != nil {
			return HandleErrorByPrintingToConsole(cmd, err)
		}

		zctx = mergeBindingKeepingFirst(zctx, dctx)

		cmd.SetContext(zctx)

		if err := callRunMethod(cmd, method, tpe); err != nil {
			return HandleErrorByPrintingToConsole(cmd, err)
		}
		return nil
	}

	if name != "" {
		cmd.Use = name
	}

	ctx = SetNamedCommand(ctx, name, &NamedCommand{
		cmd:        cmd,
		method:     method,
		methodType: tpe,
	})

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

	flagb := getFlagBindings(ctx)

	for name, cmd := range named {
		if name == RootCommandName {
			continue
		}

		for _, arg := range listOfArgs(cmd.methodType) {
			if flagb[arg] == nil {
				continue
			}
			flagb[arg](cmd.cmd.Flags())
		}

		rootcmd.cmd.AddCommand(cmd.cmd)
	}

	// we keep the context so we can use it to resolve bindings when a command is run
	rootcmd.ptr = ctx

	return rootcmd.cmd
}

func MustNewCommand(ctx context.Context, name string, snk Snakeable) context.Context {
	ctx, err := NewCommand(ctx, name, snk)
	if err != nil {
		panic(err)
	}
	return ctx
}

// type funcSnakeable struct {
// 	fn  func(...any) error
// 	cmd *cobra.Command
// }

// func (f *funcSnakeable) PreRun(ctx context.Context, args []string) (context.Context, error) {
// 	return ctx, nil
// }

// func (f *funcSnakeable) Register(ctx context.Context) (*cobra.Command, error) {
// 	return f.cmd, nil
// }

// func (f *funcSnakeable) Run(ctx context.Context, args []string) error {

// func MustNewCommandFunc(ctx context.Context, name string, cmd *cobra.Command, fn func(...any) error) context.Context {
// 	return MustNewCommand(ctx, name, &funcSnakeable{
// 		fn: fn,
// 	})
// }

func WithRootCommand(ctx context.Context, x func(*cobra.Command) error) error {
	root := GetRootCommand(ctx)
	if root == nil {
		return fmt.Errorf("snake.WithRootCommand: no root command found in context")
	}
	return x(root.cmd)
}

func WithNamedCommand(ctx context.Context, name string, x func(*cobra.Command) error) error {
	cmd := GetNamedCommand(ctx, name)
	if cmd == nil {
		return fmt.Errorf("snake.WithNamedCommand: no named command found in context")
	}
	return x(cmd.cmd)
}

func WithActiveCommand(ctx context.Context, x func(*cobra.Command) error) error {
	cmd := GetActiveNamedCommand(ctx)
	if cmd == nil {
		return fmt.Errorf("snake.WithActiveCommand: no active command found in context")
	}
	return x(cmd.cmd)
}
