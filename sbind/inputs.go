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
type StringEnumInput = enumInput[string]

func MethodName(m ValidatedRunMethod) string {
	return reflect.ValueOf(m.RunFunc()).Type().String()
}

func DependancyInputs[G ValidatedRunMethod](str string, m FMap[G], enum ...EnumOption) ([]Input, error) {
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

func InputsFor(m ValidatedRunMethod, enum ...EnumOption) ([]Input, error) {
	flds := StructFields(m)
	shared := MenthodIsShared(m)
	resp := make([]Input, 0)
	for _, f := range flds {

		if !f.IsExported() {
			continue
		}

		if f.Type.Kind() == reflect.Ptr {
			return nil, errors.Errorf("field %q in %v is a pointer type", f.Name, m)
		}

		switch f.Type.Kind() {
		case reflect.String:
			if f.Type.Name() != "string" {
				en, err := NewGenericEnumInput[string](f, m, shared, enum...)
				if err != nil {
					return nil, err
				}
				resp = append(resp, en)
			} else {
				resp = append(resp, NewSimpleValueInput[string](f, m, shared))
			}
		case reflect.Int:
			resp = append(resp, NewSimpleValueInput[int](f, m, shared))
		case reflect.Bool:
			resp = append(resp, NewSimpleValueInput[bool](f, m, shared))
		case reflect.Array, reflect.Slice:
			resp = append(resp, NewSimpleValueInput[[]string](f, m, shared))
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

type enumInput[T EnumConstraint] struct {
	simpleValueInput[T]
	options     []T
	rawTypeName string
}

func (me *enumInput[T]) Options() []T {
	return me.options
}

func getEnumOptionsFrom[T EnumConstraint](mytype reflect.Type, enum ...EnumOption) ([]T, error) {
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

	// this is not only ever used for strings, but it is the only type that is currently supported
	myTarget := reflect.TypeOf((*T)(nil)).Elem()

	res := make([]T, 0)
	for _, v := range sel.Options() {
		if !reflect.ValueOf(v).CanConvert(myTarget) {
			return nil, errors.Errorf("cannot convert type %T into %q", v, rawTypeName)
		}
		res = append(res, reflect.ValueOf(v).Convert(myTarget).Interface().(T))
	}
	return res, nil
}

func (me *enumInput[T]) ApplyOptions(opts []EnumOption) error {
	var sel EnumOption

	mytype := reflect.TypeOf(me.val).Elem()

	for _, v := range opts {
		if v.RawTypeName() != me.rawTypeName {
			continue
		}

		sel = v
	}
	if sel == nil {
		return errors.Errorf("no options for %q", me.rawTypeName)
	}

	for _, v := range sel.Options() {
		if !reflect.ValueOf(v).CanConvert(mytype) {
			return errors.Errorf("cannot convert type %T into %q", v, me.rawTypeName)
		}
		me.options = append(me.options, reflect.ValueOf(v).Convert(mytype).Interface().(T))
	}
	return nil
}

func NewGenericEnumInput[T EnumConstraint](f reflect.StructField, m ValidatedRunMethod, shared bool, enum ...EnumOption) (*enumInput[T], error) {

	v := FieldByName(m, f.Name)

	mytype := v.Type()

	opts, err := getEnumOptionsFrom[T](mytype, enum...)
	if err != nil {
		return nil, err
	}

	return &enumInput[T]{
		simpleValueInput: simpleValueInput[T]{
			genericInput: NewGenericInput(f, m, shared),
			val:          (*T)(v.Addr().UnsafePointer()),
		},
		options:     opts,
		rawTypeName: mytype.String(),
	}, nil
}

func (me *enumInput[T]) Ptr() any {
	return me.val
}

func NewSimpleValueInput[T any](f reflect.StructField, m ValidatedRunMethod, shared bool) *simpleValueInput[T] {
	v := FieldByName(m, f.Name)

	return &simpleValueInput[T]{
		genericInput: NewGenericInput(f, m, shared),
		val:          v.Addr().Interface().(*T),
	}
}

func NewGenericInput(f reflect.StructField, m ValidatedRunMethod, shared bool) *genericInput {
	return &genericInput{
		field:  f,
		parent: MethodName(m),
		shared: shared,
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

func (me *simpleValueInput[T]) Default() T {
	defstr := me.field.Tag.Get("default")
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
