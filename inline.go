package snake

import (
	"github.com/spf13/pflag"
)

var _ Flagged = (*inlineResolver[any])(nil)

type inlineResolver[I any] struct {
	flagFunc func(*pflag.FlagSet)
	runFunc  func() (I, error)
}

func (me *inlineResolver[I]) GetFlag(str string) any {
	g := &pflag.FlagSet{}
	me.flagFunc(g)
	return g.Lookup(str).Value.(any)
}

func (me *inlineResolver[I]) SetFlag(str string, val any) {
	g := &pflag.FlagSet{}
	g.AddFlag(&pflag.Flag{
		Name:  str,
		Value: val.(pflag.Value),
	})
	me.flagFunc(g)
}

func (me *inlineResolver[I]) Run() (I, error) {
	return me.runFunc()
}

func NewArgInlineFunc[I any](flagFunc func(*pflag.FlagSet), runFunc func() (I, error)) Flagged {
	return &inlineResolver[I]{flagFunc: flagFunc, runFunc: runFunc}
}
