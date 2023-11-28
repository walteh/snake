package sbind

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
