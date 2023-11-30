package sbind

import (
	"context"
	"reflect"
)

type NewSnakeOpts struct {
	Resolvers                  []Resolver
	NamedResolvers             map[string]Resolver
	GlobalContextResolverFlags bool
}

type Snake interface {
	ResolverNames() []string
	Resolve(string) Resolver
	// Bound(string) *reflect.Value
	Binder() *Binder
}

type defaultSnake struct {
	bindings  *Binder
	resolvers map[string]Resolver
}

func (me *defaultSnake) ResolverNames() []string {
	names := make([]string, 0)
	for k := range me.resolvers {
		names = append(names, k)
	}
	return names
}

func (me *defaultSnake) Resolve(name string) Resolver {
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
}

func NewSnake[M Method](opts *NewSnakeOpts, impl SnakeImplementation[M]) (Snake, error) {

	snk := &defaultSnake{
		bindings:  NewBinder(),
		resolvers: make(map[string]Resolver),
	}

	enums := make([]EnumOption, 0)

	// we always want context to get resolved first
	opts.NamedResolvers["root"] = MustGetRunMethod(NewNoopAsker[context.Context]())

	for _, runner := range opts.Resolvers {

		retrn := ListOfReturns(runner)

		// every return value marks this runner as the resolver for that type
		for _, r := range retrn {
			if r.Kind().String() == "error" {
				continue
			}
			snk.resolvers[reflectTypeString(r)] = runner
		}

		// enum options are also resolvers so they are passed here
		if mp, ok := runner.(EnumOption); ok {
			enums = append(enums, mp)
		}
	}

	for k, v := range opts.NamedResolvers {
		snk.resolvers[k] = v
	}

	for _, sexer := range snk.ResolverNames() {
		exer := snk.Resolve(sexer)

		if cmd, ok := exer.Ref().(M); ok {
			inpts, err := DependancyInputs(sexer, snk.Resolve, enums...)
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
