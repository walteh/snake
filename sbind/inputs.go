package sbind

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/go-faster/errors"
)

type Input interface {
	Name() string
	Type() reflect.Type
	Shared() bool
	M() Method
	Usage() string
	Ptr() any
}

type genericInput[T any] struct {
	field  reflect.StructField
	shared bool
	m      Method
	val    *T
}

type StringInput = genericInput[string]
type IntInput = genericInput[int]
type BoolInput = genericInput[bool]

func MethodName(m Method) string {
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
					return nil, errors.Errorf("multiple inputs named %q [%q, %q]", v.Name(), z, MethodName(v.M()))
				}
				nameReserved[v.Name()] = MethodName(v.M())
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

		switch f.Type.Kind() {
		case reflect.String:
			resp = append(resp, NewGenericInput[M, string](f, m, shared))
		case reflect.Int:
			resp = append(resp, NewGenericInput[M, int](f, m, shared))
		case reflect.Bool:
			resp = append(resp, NewGenericInput[M, bool](f, m, shared))
		default:
			return nil, errors.Errorf("field %q in %v is not a string or int", f.Name, m)
		}

	}
	return resp, nil
}

func NewGenericInput[M Method, T any](f reflect.StructField, m M, shared bool) *genericInput[T] {
	return &genericInput[T]{
		field:  f,
		m:      m,
		shared: shared,
		val:    reflect.Indirect(reflect.ValueOf(m).Elem()).FieldByName(f.Name).Addr().Interface().(*T),
	}
}

func (me *genericInput[T]) Value() *T {
	return me.val
}

func (me *genericInput[T]) Name() string {
	return strings.ToLower(me.field.Name)
}

func (me *genericInput[T]) Type() reflect.Type {
	return me.field.Type
}

func (me *genericInput[T]) Shared() bool {
	return me.shared
}

func (me *genericInput[T]) M() Method {
	return me.m
}

func (me *genericInput[T]) Usage() string {
	return me.field.Tag.Get("usage")
}

func (me *genericInput[T]) Ptr() any {
	return me.val
}

func (me *genericInput[T]) Default() T {
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
