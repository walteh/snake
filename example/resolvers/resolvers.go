package resolvers

import (
	"context"
	"io"

	"github.com/walteh/snake"
	"github.com/walteh/snake/example/root/sample"
)

func LoadResolvers() []snake.Resolver {
	return []snake.Resolver{

		// SINGLE RESOLVERS
		snake.MustGetResolverFor[context.Context](&ContextResolver{}),
		snake.MustGetResolverFor[CustomInterface](&CustomResolver{}),

		// MULTI RESOLVERS
		snake.MustGetResolverFor2[io.Reader, io.Writer](&DoubleResolver{}),
		snake.MustGetResolverFor3[io.ByteReader, io.ByteWriter, io.ByteScanner](&TripleResolver{}),

		// ENUM RESOLVERS
		snake.NewEnumOptionWithResolver(
			"sample-enum", "the sample of an enum",
			sample.SampleEnumX,
			sample.SampleEnumY,
			sample.SampleEnumZ,
		),
	}
}
