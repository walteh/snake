package sbind

import (
	"reflect"

	"github.com/go-faster/errors"
)

func ListOfArgs(typ reflect.Type) []reflect.Type {
	var args []reflect.Type

	for i := 0; i < typ.NumIn(); i++ {
		args = append(args, typ.In(i))
	}

	return args
}

func ListOfReturns(typ reflect.Type) []reflect.Type {
	var args []reflect.Type

	for i := 0; i < typ.NumOut(); i++ {
		args = append(args, typ.Out(i))
	}

	return args
}

func GetRunMethod(inter any) (reflect.Value, error) {
	prov, ok := inter.(MethodProvider)
	if ok {
		return prov.Method(), nil
	}

	value := reflect.ValueOf(inter)

	method := value.MethodByName("Run")
	if !method.IsValid() {
		if value.CanAddr() {
			method = value.Addr().MethodByName("Run")
		}
	}

	if !method.IsValid() {
		return reflect.Value{}, errors.Errorf("missing Run method on %q", value.Type())
	}

	return method, nil
}
