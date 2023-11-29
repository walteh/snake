package sbind

import (
	"reflect"

	"github.com/go-faster/errors"
)

type ValidatedRunMethod interface {
	RunFunc() reflect.Value
	Ref() Method
}

func MustGetRunMethod[M Method](inter M) TypedValidatedRunMethod[M] {
	m, err := GetRunMethod(inter)
	if err != nil {
		panic(err)
	}
	return m
}

type TypedValidatedRunMethod[M Method] interface {
	ValidatedRunMethod
	TypedRef() M
}

func GetRunMethod[M Method](inter M) (TypedValidatedRunMethod[M], error) {

	prov, ok := any(inter).(MethodProvider)
	if ok {
		return &runMethod[M]{
			runfunc: prov.Method(),
			strc:    inter,
		}, nil
	}

	value := reflect.ValueOf(inter)

	method := value.MethodByName("Run")
	if !method.IsValid() {
		if value.CanAddr() {
			method = value.Addr().MethodByName("Run")
		}
	}

	if !method.IsValid() {
		return nil, errors.Errorf("missing Run method on %q", value.Type())
	}

	return &runMethod[M]{
		runfunc: method,
		strc:    inter,
	}, nil

}

type runMethod[M Method] struct {
	runfunc reflect.Value
	strc    M
}

func (me *runMethod[M]) RunFunc() reflect.Value {
	return me.runfunc
}

func (me *runMethod[M]) Ref() Method {
	return me.strc
}

func (me *runMethod[M]) TypedRef() M {
	return me.strc
}

func ListOfArgs(m ValidatedRunMethod) []reflect.Type {
	var args []reflect.Type
	typ := m.RunFunc().Type()
	for i := 0; i < typ.NumIn(); i++ {
		args = append(args, typ.In(i))
	}

	return args
}

func ListOfReturns(m ValidatedRunMethod) []reflect.Type {
	var args []reflect.Type
	typ := m.RunFunc().Type()
	for i := 0; i < typ.NumOut(); i++ {
		args = append(args, typ.Out(i))
	}
	return args
}

func MenthodIsShared(run ValidatedRunMethod) bool {
	rets := ListOfReturns(run)
	// right now this logic relys on the fact that commands only return one value (the error)
	// and shared methods return two or more (the error and the values)
	if len(rets) == 1 ||
		// this is the logic to support the new Output type
		(len(rets) == 2 && rets[0].String() == reflect.TypeOf((*Output)(nil)).Elem().String()) {
		return false
	} else {
		return true
	}
}

func ReturnArgs(me ValidatedRunMethod) []reflect.Type {
	return ListOfReturns(me)
}

func FieldByName(me ValidatedRunMethod, name string) reflect.Value {
	return reflect.Indirect(reflect.ValueOf(me.Ref()).Elem()).FieldByName(name)
}

func CallMethod(me ValidatedRunMethod, args []reflect.Value) []reflect.Value {
	return me.RunFunc().Call(args)
}

func StructFields(me ValidatedRunMethod) []reflect.StructField {
	typ := reflect.TypeOf(me.Ref())
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	vis := reflect.VisibleFields(typ)
	return vis
}
