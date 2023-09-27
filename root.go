package snake

import (
	"context"

	"github.com/spf13/cobra"
)

type rootKeyT struct {
}

func SetRootCommand(ctx context.Context, cmd *cobra.Command) context.Context {
	return context.WithValue(ctx, &rootKeyT{}, cmd)
}

func GetRootCommand(ctx context.Context) *cobra.Command {
	p, ok := ctx.Value(&rootKeyT{}).(*cobra.Command)
	if ok {
		return p
	}
	return nil
}

type namedCommandKeyT struct {
	name string
}

func SetNamedCommand(ctx context.Context, name string, cmd *cobra.Command) context.Context {
	return context.WithValue(ctx, &namedCommandKeyT{name}, cmd)
}

func GetNamedCommand(ctx context.Context, name string) *cobra.Command {
	p, ok := ctx.Value(&namedCommandKeyT{name}).(*cobra.Command)
	if ok {
		return p
	}
	return nil
}

type activeCommandKeyT struct {
}

func SetActiveCommand(ctx context.Context, str string) context.Context {
	return context.WithValue(ctx, &activeCommandKeyT{}, str)
}

func GetActiveCommand(ctx context.Context) string {
	p, ok := ctx.Value(&activeCommandKeyT{}).(string)
	if ok {
		return p
	}
	return ""
}

func ClearActiveCommand(ctx context.Context) context.Context {
	return context.WithValue(ctx, &activeCommandKeyT{}, "")
}

func GetActiveNamedCommand(ctx context.Context) *cobra.Command {
	return GetNamedCommand(ctx, GetActiveCommand(ctx))
}
