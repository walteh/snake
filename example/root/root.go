package root

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/walteh/snake"
	"github.com/walteh/snake/example/resolvers"
	"github.com/walteh/snake/example/root/sample"
	"github.com/walteh/snake/scobra"
	"github.com/walteh/snake/smiddleware"
)

func NewCommand(ctx context.Context) (snake.Snake, *scobra.CobraSnake, *sample.Handler, error) {

	cmd := &cobra.Command{
		Use: "root",
	}

	impl := scobra.NewCobraSnake(cmd)

	handler := &sample.Handler{}

	commands := []snake.Resolver{
		scobra.NewCommandResolver(handler).WithMiddleware(smiddleware.NewIntervalMiddleware()),
	}

	snk, err := snake.NewSnakeWithOpts(ctx, impl, &snake.NewSnakeOpts{
		Resolvers: append(commands, resolvers.LoadResolvers()...),
	})
	if err != nil {
		return nil, nil, nil, err
	}

	return snk, impl, handler, err
}
