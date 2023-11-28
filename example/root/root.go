package root

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/spf13/cobra"

	"github.com/walteh/snake/example/root/sample"
	"github.com/walteh/snake/sbind"
	"github.com/walteh/snake/scobra"
)

func NewCommand(ctx context.Context) (*cobra.Command, *sample.Handler, error) {

	cmd := &cobra.Command{
		Use: "retab",
	}

	handler := &sample.Handler{}

	out, err := scobra.NewCobraSnake(cmd, &scobra.NewSCobraOpts{
		Commands: []scobra.SCobra{
			handler,
		},
		Resolvers: []sbind.Method{
			&ContextResolver{},
			&CustomResolver{},
			&DoubleResolver{},
			&TripleResolver{},
			&EnumResolver{},
		},
		EnumTypeFunc: func(s string) ([]any, error) {
			if s == "sample.SampleEnum" {
				return []any{
					sample.SampleEnumX,
					sample.SampleEnumY,
					sample.SampleEnumZ,
				}, nil
			}
			return nil, errors.Errorf("unknown enum type %q", s)

		},
	})

	return out, handler, err
}
