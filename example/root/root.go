package root

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/walteh/snake/example/resolvers"
	"github.com/walteh/snake/example/root/sample"
	"github.com/walteh/snake/sbind"
	"github.com/walteh/snake/scobra"
)

func NewCommand(ctx context.Context) (*cobra.Command, *sample.Handler, error) {

	cmd := &cobra.Command{
		Use: "root",
	}

	handler := &sample.Handler{}

	out, err := scobra.NewCobraSnake(cmd)
	if err != nil {
		return nil, nil, err
	}

	loaded := resolvers.LoadResolvers()

	commands := []sbind.Resolver{
		scobra.NewCommandResolver(handler),
	}

	commands = append(commands, loaded...)

	_, err = sbind.NewSnake(ctx, out, commands...)
	if err != nil {
		return nil, nil, err
	}

	return out.Command, handler, err
}
