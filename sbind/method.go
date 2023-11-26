package sbind

import (
	"context"
	"reflect"
)

// type Method interface {
// 	Run() reflect.Value
// 	// ValidateResponse() error
// 	// HandleResponse([]reflect.Value) ([]*reflect.Value, error)
// 	Names() []string
// 	// Command() Cobrad
// }

// var _ Method = (*method)(nil)

// func (me *method) Run() reflect.Value {
// 	return me.method
// }

// type namedMethod struct {
// 	name   string
// 	method any
// }

// type MethodWrapper interface {
// 	Method() any
// }

// func (me *namedMethod) Method() any {
// 	return me.method
// }

// func (me *namedMethod) Names() []string {
// 	return []string{me.name}
// }

// func NewNamedMethod(method any) Method {
// 	return &namedMethod{name: reflectTypeString(reflect.TypeOf(method)), method: method}
// }

func RunArgs(me Method) []reflect.Type {
	return ListOfArgs(GetRunMethod(me).Type())
}

func ReturnArgs(me Method) []reflect.Type {
	return ListOfReturns(GetRunMethod(me).Type())
}

// func (me *method) ValidateResponse() error {
// 	return me.validationStrategy(me.ReturnArgs())
// }

// func (me *method) HandleResponse(out []reflect.Value) ([]*reflect.Value, error) {
// 	output := make([]*reflect.Value, 0)
// 	for _, v := range out {
// 		output = append(output, &v)
// 	}
// 	return output, nil
// }

// func (me *method) Names() []string {
// 	return me.names
// }

func IsContextResolver(me Method) bool {
	returns := ReturnArgs(me)
	return len(returns) == 2 && returns[0] == reflect.TypeOf((*context.Context)(nil)).Elem()
}

// type method struct {
// 	names  []string
// 	method reflect.Value

// 	setFlag            func(string, any)
// 	getFlag            func(string) any
// 	validationStrategy func([]reflect.Type) error
// 	responseStrategy   func([]reflect.Value) ([]*reflect.Value, error)
// 	// cmd                Cobrad
// }
