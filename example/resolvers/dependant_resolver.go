package resolvers

import (
	"fmt"
	"strings"
)

// var _ scobra.Flagged = (*ContextResolver)(nil)

type DependantResolver struct {
	count int
}

type DependantResolverString string

func (me *DependantResolver) Run(en SampleEnum) (DependantResolverString, error) {
	me.count++
	return DependantResolverString(strings.ToUpper(string(en)) + fmt.Sprintf(" %d", me.count)), nil
}
