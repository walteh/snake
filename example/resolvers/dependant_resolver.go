package resolvers

import (
	"fmt"
	"strings"

	"github.com/walteh/snake"
)

func DependantRunner() snake.Runner {
	return snake.GenRunResolver_In01_Out02(&DependantResolver{})
}

type DependantResolver struct {
	count int
}

type DependantResolverString string

func (me *DependantResolver) Run(en SampleEnum) (DependantResolverString, error) {
	me.count++
	return DependantResolverString(strings.ToUpper(string(en)) + fmt.Sprintf(" %d", me.count)), nil
}
