package sbind

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/go-faster/errors"
)

type Input interface {
	Name() string
	Type() reflect.Type
	Shared() bool
	M() Method
	Usage() string
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

func DependancyInputs[G Method](str string, m FMap[G]) ([]Input, error) {
	deps, err := DependanciesOf(str, m)
	if err != nil {
		return nil, err
	}

	procd := make([]Input, 0)
	for _, f := range deps {
		if z := m(f); reflect.ValueOf(z).IsNil() {
			return nil, errors.Errorf("missing resolver for %q", f)
		} else {
			inp, err := InputsFor(z)
			if err != nil {
				return nil, err
			}
			for _, v := range inp {
				procd = append(procd, v)
			}
		}
	}

	return procd, nil
}

// func InputsFor(m Method) (map[string]reflect.Type, error) {
// 	flds := reflect.VisibleFields(reflect.TypeOf(m))
// 	mname := reflect.TypeOf(m).String()
// 	resp := make(map[string]reflect.Type, 0)
// 	for _, f := range flds {
// 		if resp[f.Name] != nil {
// 			return nil, errors.Errorf("duplicate field %q in %q", f.Name, mname)
// 		}
// 		resp[f.Name] = f.Type
// 	}
// 	return resp, nil
// }

func InputsFor[M Method](m M) ([]Input, error) {
	typ := reflect.TypeOf(m)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	fmt.Println("InputsFor", typ)
	flds := reflect.VisibleFields(typ)
	shared := MethodIsShared(m)
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
	return me.field.Name
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

// func InputValueAs[T any](me Input[T]) *T {
// 	return me.Value()
// }

func MethodIsShared(m Method) bool {
	run := GetRunMethod(m)
	rets := ListOfReturns(run.Type())
	// right now this logic relys on the fact that commands only return one value (the error)
	// and shared methods return two or more (the error and the values)
	if len(rets) == 1 {
		return false
	} else {
		return true
	}
}
