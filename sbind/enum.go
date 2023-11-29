package sbind

import "reflect"

type EnumConstraint interface {
	~string
}

type EnumOption interface {
	RawTypeName() string
	Options() []string
	NoopResolver() Method
}

type rawEnumOption[T EnumConstraint] struct {
	rawTypeName string
	options     []T
}

func NewEnumOption[T EnumConstraint](input ...T) EnumOption {
	return &rawEnumOption[T]{
		rawTypeName: reflect.TypeOf((*T)(nil)).Elem().String(),
		options:     input,
	}
}

func (me *rawEnumOption[T]) RawTypeName() string {
	return me.rawTypeName
}

func (me *rawEnumOption[T]) Options() []string {
	opts := make([]string, len(me.options))
	for i, v := range me.options {
		opts[i] = string(v)
	}
	return opts
}

func (me *rawEnumOption[T]) OptionsTyped() []T {
	return me.options
}

func (me *rawEnumOption[T]) NoopResolver() Method {
	return NewNoopMethod[T]()
}
