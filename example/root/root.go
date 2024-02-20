package root

import (
	"context"
	"time"

	"github.com/spf13/cobra"

	"github.com/walteh/snake"
	"github.com/walteh/snake/example/resolvers"
	"github.com/walteh/snake/example/root/basic"
	"github.com/walteh/snake/example/root/sample"
	"github.com/walteh/snake/scobra"
	"github.com/walteh/snake/smiddleware"
	"github.com/walteh/terrors"
)

func NewCommand(ctx context.Context) (snake.Snake, *scobra.CobraSnake, *sample.Handler, error) {

	cmd := &cobra.Command{
		Use: "root",
	}

	impl := scobra.NewCobraSnake(ctx, cmd)

	handler := &sample.Handler{}

	sobps := snake.Opts(
		snake.Commands(
			scobra.NewCommand(&basic.Handler{}),
			scobra.NewCommand(handler).WithMiddleware(smiddleware.NewIntervalMiddlewareWithDefault(time.Second)),
		),
		snake.Resolvers(
			resolvers.CustomRunner(),
			resolvers.TripleRunner(),
			resolvers.DoubleRunner(),
			resolvers.DependantRunner(),
			snake.NewEnumOptionWithResolver(
				"sample-enum", "the sample of an enum",
				resolvers.SampleEnumX,
				resolvers.SampleEnumY,
				resolvers.SampleEnumZ,
			),
		))

	sobps.OverrideEnumResolver = func(name string, values []string) (string, error) {
		if name == "resolvers.SampleEnum" {
			return "y", nil
		}
		return "", terrors.Errorf("unknown enum %s", name)
	}

	snk, err := snake.NewSnakeWithOpts(ctx, impl, sobps)
	if err != nil {
		return nil, nil, nil, err
	}

	return snk, impl, handler, err
}
