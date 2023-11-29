package sbind

import (
	"context"
	"reflect"
)

type EnumTypeFunc func(string) ([]any, error)

type EnumOptionResolver func(string) ([]string, error)

type NewSnakeOpts struct {
	Resolvers                  []ValidatedRunMethod
	NamedResolvers             map[string]ValidatedRunMethod
	GlobalContextResolverFlags bool
	Enums                      []EnumOption
}

type Snake interface {
	ResolverNames() []string
	Resolve(string) ValidatedRunMethod
	// Bound(string) *reflect.Value
	Binder() *Binder
	SetResolver(string, ValidatedRunMethod)
}

type defaultSnake struct {
	bindings  *Binder
	resolvers map[string]ValidatedRunMethod
}

func (me *defaultSnake) SetResolver(name string, meth ValidatedRunMethod) {
	me.resolvers[name] = meth
}

func (me *defaultSnake) ResolverNames() []string {
	names := make([]string, 0)
	for k := range me.resolvers {
		names = append(names, k)
	}
	return names
}

func (me *defaultSnake) Resolve(name string) ValidatedRunMethod {
	return me.resolvers[name]
}

func (me *defaultSnake) Binder() *Binder {
	return me.bindings
}

type MethodProvider interface {
	Method() reflect.Value
}

type SnakeImplementation[X any] interface {
	Decorate(X, Snake, []Input) error
	// ProcessInputs(X, Snake) error
}

func NewSnake[M Method](opts *NewSnakeOpts, impl SnakeImplementation[M]) (Snake, error) {

	snk := &defaultSnake{
		bindings:  NewBinder(),
		resolvers: make(map[string]ValidatedRunMethod),
	}

	// we always want context to get resolved first
	opts.NamedResolvers["root"] = MustGetRunMethod(NewNoopAsker[context.Context]())

	for _, runner := range opts.Resolvers {

		retrn := ReturnArgs(runner)

		for _, r := range retrn {
			if r.Kind().String() == "error" {
				continue
			}
			snk.resolvers[reflectTypeString(r)] = runner
		}
	}

	for k, v := range opts.NamedResolvers {
		snk.resolvers[k] = v
	}

	for _, sexer := range snk.ResolverNames() {
		exer := snk.Resolve(sexer)

		if cmd, ok := exer.Ref().(M); ok {
			inpts, err := DependancyInputs(sexer, snk.Resolve, opts.Enums...)
			if err != nil {
				return nil, err
			}

			err = impl.Decorate(cmd, snk, inpts)
			if err != nil {
				return nil, err
			}

			continue
		}

	}

	return snk, nil

}
