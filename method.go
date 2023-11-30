package snake

import (
	"reflect"

	"github.com/go-faster/errors"
)

type TypedResolver[M Method] interface {
	Resolver
	TypedRef() M
}

type Resolver interface {
	RunFunc() reflect.Value
	Ref() Method
	IsResolver()
}

func MustGetRunMethod[M Method](inter M) TypedResolver[M] {
	m, err := getRunMethod(inter)
	if err != nil {
		panic(err)
	}
	return m
}

func MustGetResolverFor[M any](inter Method) Resolver {
	return mustGetResolverForRaw(inter, (*M)(nil))
}

func MustGetResolverFor2[M1, M2 any](inter Method) Resolver {
	return mustGetResolverForRaw(inter, (*M1)(nil), (*M2)(nil))
}

func MustGetResolverFor3[M1, M2, M3 any](inter Method) Resolver {
	return mustGetResolverForRaw(inter, (*M1)(nil), (*M2)(nil), (*M3)(nil))
}

func mustGetResolverForRaw(inter any, args ...any) Resolver {
	run, err := getRunMethod(inter)
	if err != nil {
		panic(err)
	}

	resvf := ResolverFor(run)

	for _, arg := range args {
		argptr := reflect.TypeOf(arg).Elem()
		if yes, ok := resvf[argptr.String()]; !ok || !yes {
			panic(errors.Errorf("%q is not a resolver for %q", reflect.TypeOf(inter).String(), argptr.String()))
		}
	}

	return run
}
func getRunMethod[M Method](inter M) (*simpleResolver[M], error) {

	prov, ok := any(inter).(MethodProvider)
	if ok {
		return &simpleResolver[M]{
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

	return &simpleResolver[M]{
		runfunc: method,
		strc:    inter,
	}, nil
}

type simpleResolver[M Method] struct {
	runfunc reflect.Value
	strc    M
}

func (me *simpleResolver[M]) RunFunc() reflect.Value {
	return me.runfunc
}

func (me *simpleResolver[M]) Ref() Method {
	return me.strc
}

func (me *simpleResolver[M]) TypedRef() M {
	return me.strc
}

func (me *simpleResolver[M]) IsResolver() {}

func ListOfArgs(m Resolver) []reflect.Type {
	var args []reflect.Type
	typ := m.RunFunc().Type()
	for i := 0; i < typ.NumIn(); i++ {
		args = append(args, typ.In(i))
	}

	return args
}

func ListOfReturns(m Resolver) []reflect.Type {
	var args []reflect.Type
	typ := m.RunFunc().Type()
	for i := 0; i < typ.NumOut(); i++ {
		args = append(args, typ.Out(i))
	}
	return args
}

func MenthodIsShared(run Resolver) bool {
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

func ResolverFor(m Resolver) map[string]bool {
	resp := make(map[string]bool, 0)
	for _, f := range ListOfReturns(m) {
		if f.String() == "error" {
			continue
		}
		resp[f.String()] = true
	}
	return resp
}

func FieldByName(me Resolver, name string) reflect.Value {
	return reflect.Indirect(reflect.ValueOf(me.Ref()).Elem()).FieldByName(name)
}

func CallMethod(me Resolver, args []reflect.Value) []reflect.Value {
	return me.RunFunc().Call(args)
}

func StructFields(me Resolver) []reflect.StructField {
	typ := reflect.TypeOf(me.Ref())
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	vis := reflect.VisibleFields(typ)
	return vis
}
