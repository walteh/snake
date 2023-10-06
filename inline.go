package snake

import "github.com/spf13/pflag"

var _ Flagged = (*inlineResolver[any])(nil)

type inlineResolver[I any] struct {
	flagFunc func(*pflag.FlagSet)
	runFunc  func() (I, error)
}

func (me *inlineResolver[I]) Flags(flgs *pflag.FlagSet) {
	me.flagFunc(flgs)
}

func (me *inlineResolver[I]) Run() (I, error) {
	return me.runFunc()
}

func NewInlineFunc[I any](flagFunc func(*pflag.FlagSet), runFunc func() (I, error)) Flagged {
	return &inlineResolver[I]{flagFunc: flagFunc, runFunc: runFunc}
}

func NewInlineFuncSimple[I any](runFunc func() (I, error)) Flagged {
	return &inlineResolver[I]{flagFunc: func(*pflag.FlagSet) {}, runFunc: runFunc}
}

func NewInlineSimple[I any](value I) Flagged {
	return &inlineResolver[I]{flagFunc: func(*pflag.FlagSet) {}, runFunc: func() (I, error) { return value, nil }}
}
