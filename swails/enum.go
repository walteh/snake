package swails

import (
	"reflect"
)

type enum struct {
	GOName string
	TSName string
}

func ReflectValueToEnum(v reflect.Type) enum {
	return enum{
		GOName: v.String(),
		TSName: v.String(),
	}
}

// func InputTypesAsEnum() any {
// 	arr := []enum{}

// 	inp := snake.AllInputTypes()

// 	for _, v := range inp {
// 		arr = append(arr, ReflectValueToEnum(v))
// 	}

// 	return arr
// }
