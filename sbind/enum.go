package sbind

import (
	"reflect"
	"slices"
	"strings"

	"github.com/go-faster/errors"
)

var (
	_ Resolver = (*rawEnumOption[string])(nil)
)

type EnumConstraint interface {
	~string
}
type EnumResolver func(string, []string) (string, error)

type EnumOption interface {
	// Resolver
	Resolver
	SetCurrent(string) error
	CurrentPtr() *string
	RawTypeName() string
	Options() []string
	Ptr() any
	DisplayName() string
}

type rawEnumOption[T EnumConstraint] struct {
	rawTypeName  string
	options      []T
	enumResolver EnumResolver
	name         string
	// Val needs to be exported value so it is picked up in inputs.go reflection logic
	Val *T
}

// Ref implements ValidatedRunMethod.
func (me *rawEnumOption[T]) Ref() Method {
	return me
}

// RunFunc implements ValidatedRunMethod.
func (me *rawEnumOption[T]) RunFunc() reflect.Value {
	return reflect.ValueOf(me.Run)
}

func NewEnumOptionWithResolver[T EnumConstraint](name string, resolver EnumResolver, input ...T) EnumOption {
	sel := new(T)
	if resolver != nil {
		*sel = T("select")
	}

	return &rawEnumOption[T]{
		rawTypeName:  reflect.TypeOf((*T)(nil)).Elem().String(),
		options:      input,
		enumResolver: resolver,
		name:         name,
		Val:          sel,
	}
}

func (me *rawEnumOption[T]) DisplayName() string {
	return me.name
}

func (me *rawEnumOption[T]) RawTypeName() string {
	return me.rawTypeName
}

func (me *rawEnumOption[T]) OptionsWithSelect() []string {
	opts := me.Options()
	if me.enumResolver != nil {
		opts = append(opts, "select")
	}
	return opts
}

func (me *rawEnumOption[T]) Options() []string {
	opts := make([]string, len(me.options))
	for i, v := range me.options {
		opts[i] = string(v)
	}
	return opts
}

func (e *rawEnumOption[I]) SetCurrent(vt string) error {
	if slices.Contains(e.OptionsWithSelect(), string(vt)) {
		*e.Val = I(vt)
		return nil
	}
	return errors.Errorf("invalid value %q, expected one of [\"%s\"]", vt, strings.Join(e.OptionsWithSelect(), "\", \""))
}

func (e *rawEnumOption[I]) CurrentPtr() *string {
	return (*string)(reflect.ValueOf(e.Val).UnsafePointer())
}

func (me *rawEnumOption[T]) Run() (T, error) {
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

func EnumOptionAsInput(me EnumOption, m *genericInput) *enumInput {
	return &enumInput{
		EnumOption:   me,
		genericInput: m,
	}
}

func (me *rawEnumOption[I]) Ptr() any {
	return me.CurrentPtr()
}

func (me *rawEnumOption[I]) IsResolver() {}
