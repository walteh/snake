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

// type StringArrayInput = EnumInput[string]
// type IntArrayInput = EnumInput[int]
type StringEnumInput = enumInput[string]
type IntEnumInput = enumInput[int]

func MethodName(m any) string {
	return reflect.ValueOf(m).Elem().Type().String()
}

func DependancyInputs[G Method](str string, m FMap[G]) ([]Input, error) {
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
			inp, err := InputsFor(z)
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

func InputsFor[M Method](m M) ([]Input, error) {
	typ := reflect.TypeOf(m)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	flds := reflect.VisibleFields(typ)
	shared, err := MethodIsShared(m)
	if err != nil {
		return nil, err
	}
	resp := make([]Input, 0)
	for _, f := range flds {

		if f.Type.Kind() == reflect.Ptr {
			return nil, errors.Errorf("field %q in %v is a pointer type", f.Name, m)
		}

		if !f.IsExported() {
			continue
		}

		switch f.Type.Kind() {
		case reflect.String:
			if f.Type.Name() != "string" {
				resp = append(resp, NewGenericEnumInput[M, string](f, m, shared))
			} else {
				resp = append(resp, NewSimpleValueInput[M, string](f, m, shared))
			}
		case reflect.Int:
			if f.Type.Name() != "int" {
				resp = append(resp, NewGenericEnumInput[M, int](f, m, shared))
			} else {
				resp = append(resp, NewSimpleValueInput[M, int](f, m, shared))
			}
		case reflect.Bool:
			resp = append(resp, NewSimpleValueInput[M, bool](f, m, shared))
		case reflect.Array, reflect.Slice:
			resp = append(resp, NewSimpleValueInput[M, []string](f, m, shared))
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

type enumInput[T any] struct {
	simpleValueInput[T]
	options     []T
	rawTypeName string
}

func (me *enumInput[T]) Options() []T {
	return me.options
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

func NewGenericEnumInput[M Method, T any](f reflect.StructField, m M, shared bool) *enumInput[T] {
	v := reflect.Indirect(reflect.ValueOf(m).Elem()).FieldByName(f.Name)
	return &enumInput[T]{
		simpleValueInput: simpleValueInput[T]{
			genericInput: NewGenericInput(f, m, shared),
			val:          (*T)(v.Addr().UnsafePointer()),
		},
		options:     make([]T, 0),
		rawTypeName: v.Type().String(),
	}
}

func (me *enumInput[T]) Ptr() any {
	return me.val
}

func NewSimpleValueInput[M Method, T any](f reflect.StructField, m M, shared bool) *simpleValueInput[T] {
	v := reflect.Indirect(reflect.ValueOf(m).Elem()).FieldByName(f.Name)

	return &simpleValueInput[T]{
		genericInput: NewGenericInput(f, m, shared),
		val:          v.Addr().Interface().(*T),
	}
}

func NewGenericInput[M Method](f reflect.StructField, m M, shared bool) *genericInput {
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

func MethodIsShared(m Method) (bool, error) {
	run, err := GetRunMethod(m)
	if err != nil {
		return false, err
	}

	rets := ListOfReturns(run.Type())
	// right now this logic relys on the fact that commands only return one value (the error)
	// and shared methods return two or more (the error and the values)
	if len(rets) == 1 {
		return false, nil
	} else {
		return true, nil
	}
}

type EnumOption interface {
	RawTypeName() string
	Options() []any
}

type rawEnumOption[T any] struct {
	rawTypeName string
	options     []T
}

func NewEnumOption[T ~string | ~int](input ...T) EnumOption {
	return &rawEnumOption[T]{
		rawTypeName: reflect.TypeOf(input[0]).String(),
		options:     input,
	}
}

func (me *rawEnumOption[T]) RawTypeName() string {
	return me.rawTypeName
}

func (me *rawEnumOption[T]) Options() []any {
	opts := make([]any, len(me.options))
	for i, v := range me.options {
		opts[i] = v
	}
	return opts
}
