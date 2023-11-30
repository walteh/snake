package snake

import (
	"reflect"
	"slices"
	"strings"

	"github.com/go-faster/errors"
)

// REFRESHABLE RESOLVER
var (
	_ Resolver = (*rawEnum[string])(nil)
)

type EnumResolverFunc func(string, []string) (string, error)

type Enum interface {
	Resolver
	SetCurrent(string) error
	CurrentPtr() *string
	RawTypeName() string
	Options() []string
	Ptr() any
	DisplayName() string
}

type rawEnum[T ~string] struct {
	rawTypeName  string
	options      []T
	enumResolver EnumResolverFunc
	name         string
	// Val needs to be exported value so it is picked up in inputs.go reflection logic
	Val *T
}

// Ref implements ValidatedRunMethod.
func (me *rawEnum[T]) Ref() Method {
	return me
}

// RunFunc implements ValidatedRunMethod.
func (me *rawEnum[T]) RunFunc() reflect.Value {
	return reflect.ValueOf(me.Run)
}

func NewEnumOptionWithResolver[T ~string](name string, resolver EnumResolverFunc, input ...T) Enum {
	sel := new(T)
	if resolver != nil {
		// this sets the default to "select" and not nil
		*sel = T("select")
	}

	return &rawEnum[T]{
		rawTypeName:  reflect.TypeOf((*T)(nil)).Elem().String(),
		options:      input,
		enumResolver: resolver,
		name:         name,
		Val:          sel,
	}
}

func (me *rawEnum[T]) DisplayName() string {
	return me.name
}

func (me *rawEnum[T]) RawTypeName() string {
	return me.rawTypeName
}

func (me *rawEnum[T]) OptionsWithSelect() []string {
	opts := me.Options()
	if me.enumResolver != nil {
		opts = append(opts, "select")
	}
	return opts
}

func (me *rawEnum[T]) Options() []string {
	opts := make([]string, len(me.options))
	for i, v := range me.options {
		opts[i] = string(v)
	}
	return opts
}

func (e *rawEnum[I]) SetCurrent(vt string) error {
	if slices.Contains(e.OptionsWithSelect(), string(vt)) {
		*e.Val = I(vt)
		return nil
	}
	return errors.Errorf("invalid value %q, expected one of [\"%s\"]", vt, strings.Join(e.OptionsWithSelect(), "\", \""))
}

func (e *rawEnum[I]) CurrentPtr() *string {
	return (*string)(reflect.ValueOf(e.Val).UnsafePointer())
}

func (me *rawEnum[T]) Run() (T, error) {
	if me.Val == nil || reflect.ValueOf(me.Val).IsNil() || *me.Val == "select" {
		if me.enumResolver == nil {
			return "", errors.Errorf("no enum resolver for %q", me.rawTypeName)
		}

		resolve, err := me.enumResolver(me.rawTypeName, me.Options())
		if err != nil {
			return "", err
		}

		if err := me.SetCurrent(resolve); err != nil {
			return "", err
		}
	}
	return *me.Val, nil
}

func EnumAsInput(me Enum, m *genericInput) *enumInput {
	return &enumInput{
		Enum:         me,
		genericInput: m,
	}
}

func (me *rawEnum[I]) Ptr() any {
	return me.CurrentPtr()
}

func (me *rawEnum[I]) IsResolver() {}
