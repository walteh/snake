package sbind

import "reflect"

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

func GetRunMethod(inter any) reflect.Value {
	value := reflect.ValueOf(inter)
	method := value.MethodByName("Run")
	if !method.IsValid() {
		if value.CanAddr() {
			method = value.Addr().MethodByName("Run")
		}
	}

	return method
}
