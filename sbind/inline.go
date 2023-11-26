package sbind

// type inlineResolver[A any] struct {
// 	flagFunc func(*pflag.FlagSet)
// 	runFunc  func() (A, error)
// }

// func (me *inlineResolver[A]) Flags() *pflag.FlagSet {
// 	flgs := pflag.NewFlagSet("inline", pflag.ContinueOnError)
// 	me.flagFunc(flgs)
// 	return flgs
// }

// func (me *inlineResolver[A]) Run() (A, error) {
// 	return me.runFunc()
// }

// func NewArgInlineFunc[A any](flagFunc func(*pflag.FlagSet), runFunc func() (A, error)) Flagged {
// 	return &inlineResolver[A]{flagFunc: flagFunc, runFunc: runFunc}
// }

type noopResolver[A any] struct {
}

// func (me *noopResolver[A]) Flags() *pflag.FlagSet {
// 	return pflag.NewFlagSet("noop", pflag.ContinueOnError)
// }

func (me *noopResolver[A]) Names() []string {
	return []string{}
}

func (me *noopResolver[A]) Run() (a A, err error) {
	return a, err
}

// func NewArgNoop[A any]() Flagged {
// 	return &noopResolver[A]{}
// }

func NewNoopMethod[A any]() Method {
	return &noopResolver[A]{}
}

type noopAsker[A any] struct {
}

func (me *noopAsker[A]) Names() []string {
	return []string{}
}

func (me *noopAsker[A]) Run(a A) (err error) {
	return err
}

func NewNoopAsker[A any]() Method {
	return &noopAsker[A]{}
}
