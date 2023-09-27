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
