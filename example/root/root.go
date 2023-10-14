package root

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/walteh/snake"
)

func NewCommand(ctx context.Context) (*cobra.Command, error) {

	cmd := &cobra.Command{
		Use: "retab",
	}

	return snake.NewSnake(&snake.NewSnakeOpts{
		Root:     cmd,
		Commands: []snake.Method{},
		Resolvers: []snake.Method{
			snake.NewArgumentMethod[context.Context](&ContextResolver{}),
			snake.NewArgumentMethod[CustomInterface](&CustomResolver{}),
		},
	})
}
