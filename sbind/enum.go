package sbind

import (
	"reflect"
	"slices"
	"strings"

	"github.com/go-faster/errors"
)

var (
	_ ValidatedRunMethod = (*rawEnumOption[string])(nil)
)

type EnumConstraint interface {
	~string
}
type EnumResolver func(string, []string) (string, error)

type EnumOption interface {
	ValidatedRunMethod
	SetCurrent(string) error
	CurrentPtr() *string
	RawTypeName() string
	Options() []string
	Ptr() any
}

type rawEnumOption[T EnumConstraint] struct {
	MyEnum       *T
	rawTypeName  string
	options      []T
	enumResolver EnumResolver
}

// Ref implements ValidatedRunMethod.
func (me *rawEnumOption[T]) Ref() Method {
	return me
}

// RunFunc implements ValidatedRunMethod.
func (me *rawEnumOption[T]) RunFunc() reflect.Value {
	return reflect.ValueOf(me.Run)
}

func NewEnumOptionWithResolver[T EnumConstraint](resolver EnumResolver, input ...T) EnumOption {
	sel := new(T)
	if resolver != nil {
		*sel = T("select")
	}

	return &rawEnumOption[T]{
		MyEnum:       sel,
		rawTypeName:  reflect.TypeOf((*T)(nil)).Elem().String(),
		options:      input,
		enumResolver: resolver,
	}
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
		*e.MyEnum = I(vt)
		return nil
	}
	return errors.Errorf("invalid value %q, expected one of [\"%s\"]", vt, strings.Join(e.OptionsWithSelect(), "\", \""))
}

func (e *rawEnumOption[I]) CurrentPtr() *string {
	return (*string)(reflect.ValueOf(e.MyEnum).UnsafePointer())
}

func (me *rawEnumOption[T]) Run() (T, error) {
	if me.MyEnum == nil || reflect.ValueOf(me.MyEnum).IsNil() || *me.MyEnum == "select" {
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
	return *me.MyEnum, nil
}

func EnumOptionAsInput(me EnumOption, gen *genericInput) *enumInput {
	return &enumInput{
		EnumOption:   me,
		genericInput: gen,
	}
}

func (me *rawEnumOption[I]) Ptr() any {
	return me.CurrentPtr()
}
