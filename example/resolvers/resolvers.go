package resolvers

import (
	"github.com/walteh/snake"
)

func LoadResolvers() []snake.UntypedResolver {
	return []snake.UntypedResolver{
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
