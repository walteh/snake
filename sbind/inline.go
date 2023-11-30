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

func (me *noopResolver[A]) TypedRef() *noopResolver[A] {
	return me
}

func NewNoopMethod[A any]() Resolver {
	return &noopResolver[A]{}
}

func (me *noopResolver[A]) IsResolver() {}

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

func NewNoopAsker[A any]() Resolver {
	return &noopAsker[A]{}
}

func (me *noopAsker[A]) IsResolver() {}

func (me *noopAsker[A]) TypedRef() *noopAsker[A] {
	return me
}
