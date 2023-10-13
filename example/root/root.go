package root

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/walteh/snake/example/root/sample"

	"github.com/walteh/snake"
)

func NewCommand(ctx context.Context) (*cobra.Command, error) {

	cmd := &cobra.Command{
		Use: "retab",
	}

	ctxd := snake.NewCtx()

	snake.NewCmdContext(ctxd, &sample.Handler{})

	snake.NewArgContext[context.Context](ctxd, &ContextResolver{})
	snake.NewArgContext[CustomInterface](ctxd, &CustomResolver{})

	err := snake.ApplyCtx(ctx, ctxd, cmd)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}
