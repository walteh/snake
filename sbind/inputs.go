package sbind

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/go-faster/errors"
)

type Input interface {
	Name() string
	Shared() bool
	Ptr() any
	Parent() string
}

type InputWithOptions interface {
	Options() []string
}

type StringInput = simpleValueInput[string]
type IntInput = simpleValueInput[int]
type BoolInput = simpleValueInput[bool]
type StringArrayInput = simpleValueInput[[]string]
type IntArrayInput = simpleValueInput[[]int]
type StringEnumInput = enumInput

func MethodName(m Resolver) string {
	return reflect.ValueOf(m.Ref()).Type().String()
}

func DependancyInputs(str string, m FMap, enum ...EnumOption) ([]Input, error) {
	deps, err := DependanciesOf(str, m)
	if err != nil {
		return nil, err
	}

	procd := make(map[any]Input, 0)
	nameReserved := make(map[string]string, 0)
	for _, f := range deps {
		if z := m(f); reflect.ValueOf(z).IsNil() {
			return nil, errors.Errorf("missing resolver for %q", f)
		} else {
			inp, err := InputsFor(z, enum...)
			if err != nil {
				return nil, err
			}
			for _, v := range inp {
				// if they are references to same value, then no need to worry about potential conflicts
				if _, ok := procd[v.Ptr()]; ok {
					continue
				}
				procd[v.Ptr()] = v
				if z, ok := nameReserved[v.Name()]; ok {
					return nil, errors.Errorf("multiple inputs named %q [%q, %q]", v.Name(), z, v.Parent())
				}
				nameReserved[v.Name()] = v.Parent()
			}
		}
	}

	resp := make([]Input, 0)
	for _, v := range procd {
		resp = append(resp, v)
	}

	return resp, nil
}

func InputsFor(m Resolver, enum ...EnumOption) ([]Input, error) {
	resp := make([]Input, 0)
	for _, f := range StructFields(m) {

		if !f.IsExported() {
			continue
		}

		if f.Type.Kind() == reflect.Ptr {
			f.Type = f.Type.Elem()
			// return nil, errors.Errorf("field %q in %T is a pointer type", f.Name, m)
		}

		switch f.Type.Kind() {
		case reflect.String:
			if f.Type.Name() != "string" {
				en, err := NewGenericEnumInput(f, m, enum...)
				if err != nil {
					return nil, err
				}
				resp = append(resp, en)
			} else {
				resp = append(resp, NewSimpleValueInput[string](f, m))
			}
		case reflect.Int:
			resp = append(resp, NewSimpleValueInput[int](f, m))
		case reflect.Bool:
			resp = append(resp, NewSimpleValueInput[bool](f, m))
		case reflect.Array, reflect.Slice:
			resp = append(resp, NewSimpleValueInput[[]string](f, m))
		default:
			return nil, errors.Errorf("field %q in %v is not a string or int", f.Name, m)
		}

	}
	return resp, nil
}

type genericInput struct {
	field  reflect.StructField
	shared bool
	parent string
}

type simpleValueInput[T any] struct {
	*genericInput
	val *T
}

type enumInput struct {
	EnumOption
	*genericInput
}

func (me *enumInput) Name() string {
	return me.EnumOption.DisplayName()
}

func getEnumOptionsFrom(mytype reflect.Type, enum ...EnumOption) (EnumOption, error) {
	rawTypeName := mytype.String()
	var sel EnumOption
	for _, v := range enum {
		if v.RawTypeName() != rawTypeName {
			continue
		}

		sel = v
	}
	if sel == nil {
		return nil, errors.Errorf("no options for %q", rawTypeName)
	}

	return sel, nil

}

func NewGenericEnumInput(f reflect.StructField, m Resolver, enum ...EnumOption) (*enumInput, error) {

	mytype := FieldByName(m, f.Name).Type()

	if mytype.Kind() == reflect.Ptr {
		mytype = mytype.Elem()
	}

	opts, err := getEnumOptionsFrom(mytype, enum...)
	if err != nil {
		return nil, err
	}

	return EnumOptionAsInput(opts, NewGenericInput(f, m)), nil
}

func NewSimpleValueInput[T any](f reflect.StructField, m Resolver) *simpleValueInput[T] {
	v := FieldByName(m, f.Name)

	return &simpleValueInput[T]{
		genericInput: NewGenericInput(f, m),
		val:          v.Addr().Interface().(*T),
	}
}

func NewGenericInput(f reflect.StructField, m Resolver) *genericInput {
	return &genericInput{
		field:  f,
		parent: MethodName(m),
		shared: MenthodIsShared(m),
	}
}

func (me *simpleValueInput[T]) Value() *T {
	return me.val
}

func (me *genericInput) Name() string {
	return strings.ToLower(me.field.Name)
}

func (me *genericInput) Shared() bool {
	return me.shared
}

func (me *genericInput) Parent() string {
	return me.parent
}

func (me *genericInput) Usage() string {
	return me.field.Tag.Get("usage")
}

func (me *simpleValueInput[T]) Ptr() any {
	return me.val
}

func (me *genericInput) Default() string {
	return me.field.Tag.Get("default")
}

func (me *simpleValueInput[T]) Default() T {
	defstr := me.genericInput.Default()
	switch any(me.val).(type) {
	case *string:
		return any(defstr).(T)
	case *int:
		if defstr == "" {
			return any(0).(T)
		}
		intt, err := strconv.Atoi(defstr)
		if err != nil {
			panic(err)
		}
		return any(intt).(T)
	case *bool:
		if defstr == "" {
			return any(false).(T)
		}
		boolt, err := strconv.ParseBool(defstr)
		if err != nil {
			panic(err)
		}
		return any(boolt).(T)
	default:
		panic("unhandled type")
	}
}