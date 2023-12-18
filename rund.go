package snake

import "reflect"

type Runner interface {
	isRunner()
	Resolver
}

type TypedRunner[X any] interface {
	Runner
	TypedRef() X
}

type NamedRunner interface {
	Runner
	Name() string
	Description() string
}

var _ Runner = (*rund[Method])(nil)

type rund[X any] struct {
	internal X
}

func (r *rund[X]) IsResolver() {}

func (r *rund[X]) isRunner() {}

func (r *rund[X]) RunFunc() reflect.Value {
	return reflect.ValueOf(r.internal).MethodByName("Run")
}

func (r *rund[X]) Ref() Method {
	return any(r.internal).(Method)
}

func (r *rund[X]) TypedRef() X {
	return r.internal
}

type Runnable[M any] interface {
	NamedMethod
	RunMethod() TypedRunner[M]
}
