package root

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/walteh/snake"
	"github.com/walteh/snake/example/root/sample"
)

func NewCommand(ctx context.Context) (*cobra.Command, *sample.Handler, error) {

	cmd := &cobra.Command{
		Use: "retab",
	}

	handler := &sample.Handler{}

	out, err := snake.NewSnake(&snake.NewSnakeOpts{
		Root: cmd,
		Commands: []snake.Method{
			snake.NewCommandMethod(handler),
		},
		Resolvers: []snake.Method{
			snake.NewArgumentMethod[context.Context](&ContextResolver{}),
			snake.NewArgumentMethod[CustomInterface](&CustomResolver{}),
			snake.New2ArgumentsMethod[io.Reader, io.Writer](&DoubleResolver{}),
		},
	})

	return out, handler, err
}
