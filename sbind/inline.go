package sbind

import "reflect"

type noopResolver[A any] struct {
}

func (me *noopResolver[A]) Names() []string {
	return []string{}
}

func (me *noopResolver[A]) Run() (a A, err error) {
	return a, err
}

func (me *noopResolver[A]) RunFunc() reflect.Value {
	return reflect.ValueOf(me.Run)
}

func (me *noopResolver[A]) Ref() Method {
	return me
}

func NewNoopMethod[A any]() ValidatedRunMethod {
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

func (me *noopAsker[A]) RunFunc() reflect.Value {
	return reflect.ValueOf(me.Run)
}

func (me *noopAsker[A]) Ref() Method {
	return me
}

func NewNoopAsker[A any]() ValidatedRunMethod {
	return &noopAsker[A]{}
}
