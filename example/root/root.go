package root

import (
	"context"
	"io"

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
		Resolvers: []sbind.Resolver{
			sbind.MustGetResolverFor[context.Context](&ContextResolver{}),
			sbind.MustGetResolverFor[CustomInterface](&CustomResolver{}),
			sbind.MustGetResolverFor2[io.Reader, io.Writer](&DoubleResolver{}),
			sbind.MustGetResolverFor3[io.ByteReader, io.ByteWriter, io.ByteScanner](&TripleResolver{}),
			sbind.NewEnumOptionWithResolver(
				"the-cool-enum",
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
