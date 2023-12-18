package resolvers

import (
	"github.com/walteh/snake"
)

func LoadResolvers() []snake.Resolver {
	return []snake.Resolver{
		DependantRunner(),
		DoubleRunner(),
		CustomRunner(),
		TripleRunner(),

		// ENUM RESOLVERS
		snake.NewEnumOptionWithResolver(
			"sample-enum", "the sample of an enum",
			SampleEnumX,
			SampleEnumY,
			SampleEnumZ,
		),
	}
}
