package sbind

import (
	"context"
	"reflect"
)

type EnumTypeFunc func(string) ([]any, error)

type EnumOptionResolver func(string) ([]string, error)

type NewSnakeOpts struct {
	Resolvers                  []Method
	NamedResolvers             map[string]Method
	GlobalContextResolverFlags bool
	Enums                      []EnumOption
}

type Snake interface {
	ResolverNames() []string
	Resolve(string) Method
	// Bound(string) *reflect.Value
	Binder() *Binder
}

type defaultSnake struct {
	bindings  *Binder
	resolvers map[string]Method
}

func (me *defaultSnake) ResolverNames() []string {
	names := make([]string, 0)
	for k := range me.resolvers {
		names = append(names, k)
	}
	return names
}

func (me *defaultSnake) Resolve(name string) Method {
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
		resolvers: make(map[string]Method),
	}

	// we always want context to get resolved first
	opts.NamedResolvers["root"] = NewNoopAsker[context.Context]()

	for _, v := range opts.Resolvers {

		retrn, err := ReturnArgs(v)
		if err != nil {
			return nil, err
		}

		for _, r := range retrn {
			if r.Kind().String() == "error" {
				continue
			}
			snk.resolvers[reflectTypeString(r)] = v
		}
	}

	for k, v := range opts.NamedResolvers {
		snk.resolvers[k] = v
	}

	for _, sexer := range snk.ResolverNames() {
		exer := snk.Resolve(sexer)

		if cmd, ok := exer.(M); ok {
			inpts, err := DependancyInputs(sexer, snk.Resolve)
			if err != nil {
				return nil, err
			}

			for _, v := range inpts {
				switch vok := v.(type) {
				case *enumInput[string]:
					err = vok.ApplyOptions(opts.Enums)
				case *enumInput[int]:
					err = vok.ApplyOptions(opts.Enums)
				}
				if err != nil {
					return nil, err
				}
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
