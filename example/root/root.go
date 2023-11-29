package root

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/walteh/snake/example/root/sample"
	"github.com/walteh/snake/sbind"
	"github.com/walteh/snake/scobra"
)

func NewCommand(ctx context.Context) (*cobra.Command, *sample.Handler, error) {

	cmd := &cobra.Command{
		Use: "root",
	}

	handler := &sample.Handler{}

	out, err := scobra.NewCobraSnake(cmd, &scobra.NewSCobraOpts{
		Commands: []scobra.SCobra{
			handler,
		},
		Resolvers: []sbind.ValidatedRunMethod{
			sbind.MustGetRunMethod(&ContextResolver{}),
			sbind.MustGetRunMethod(&CustomResolver{}),
			sbind.MustGetRunMethod(&DoubleResolver{}),
			sbind.MustGetRunMethod(&TripleResolver{}),
		},
		Enums: []sbind.EnumOption{
			sbind.NewEnumOptionWithResolver(
				func(s1 string, s2 []string) (string, error) {
					return string(sample.SampleEnumY), nil
				},
				sample.SampleEnumX,
				sample.SampleEnumY,
				sample.SampleEnumZ,
			),
		},
	})

	return out, handler, err
}
