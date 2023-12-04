package resolvers

import (
	"io"

	"github.com/walteh/snake"
)

func LoadResolvers() []snake.Resolver {
	return []snake.Resolver{

		// SINGLE RESOLVERS
		// snake.MustGetResolverFor[context.Context](&ContextResolver{}),
		snake.MustGetResolverFor[DependantResolverString](&DependantResolver{}),
		snake.MustGetResolverFor[CustomInterface](&CustomResolver{}),

		// MULTI RESOLVERS
		snake.MustGetResolverFor2[io.Reader, io.Writer](&DoubleResolver{}),
		snake.MustGetResolverFor3[io.ByteReader, io.ByteWriter, io.ByteScanner](&TripleResolver{}),

		// ENUM RESOLVERS
		snake.NewEnumOptionWithResolver(
			"sample-enum", "the sample of an enum",
			SampleEnumX,
			SampleEnumY,
			SampleEnumZ,
		),
	}
}
