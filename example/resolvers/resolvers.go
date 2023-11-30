package resolvers

import (
	"context"
	"io"

	"github.com/walteh/snake/example/root/sample"
	"github.com/walteh/snake/sbind"
)

func LoadResolvers() []sbind.Resolver {
	return []sbind.Resolver{

		// SINGLE RESOLVERS
		sbind.MustGetResolverFor[context.Context](&ContextResolver{}),
		sbind.MustGetResolverFor[CustomInterface](&CustomResolver{}),

		// MULTI RESOLVERS
		sbind.MustGetResolverFor2[io.Reader, io.Writer](&DoubleResolver{}),
		sbind.MustGetResolverFor3[io.ByteReader, io.ByteWriter, io.ByteScanner](&TripleResolver{}),

		// ENUM RESOLVERS
		sbind.NewEnumOptionWithResolver(
			"sample-enum",
			func(s1 string, s2 []string) (string, error) {
				return string(sample.SampleEnumY), nil
			},
			sample.SampleEnumX,
			sample.SampleEnumY,
			sample.SampleEnumZ,
		),
	}
}
