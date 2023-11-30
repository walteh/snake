package sbind

import (
	"context"
	"reflect"
)

type NewSnakeOpts[M NamedMethod] struct {
	Resolvers      []Resolver
	Implementation SnakeImplementation[M]
}

type Snake interface {
	ResolverNames() []string
	Resolve(string) Resolver
}

type defaultSnake struct {
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

type MethodProvider interface {
	Method() reflect.Value
}

type SnakeImplementation[X any] interface {
	Decorate(context.Context, X, Snake, []Input) error
	ManagedResolvers(context.Context) []Resolver
	OnSnakeInit(context.Context, Snake) error
}

func NewSnake[M NamedMethod](ctx context.Context, impl SnakeImplementation[M], res ...Resolver) (Snake, error) {
	return NewSnakeWithOpts(ctx, &NewSnakeOpts[M]{
		Resolvers:      res,
		Implementation: impl,
	})
}

func NewSnakeWithOpts[M NamedMethod](ctx context.Context, opts *NewSnakeOpts[M]) (Snake, error) {

	snk := &defaultSnake{
		resolvers: make(map[string]Resolver),
	}

	enums := make([]EnumOption, 0)

	named := make(map[string]TypedResolver[M])

	inputResolvers := append(opts.Resolvers, opts.Implementation.ManagedResolvers(ctx)...)

	for _, runner := range inputResolvers {

		if nmd, err := runner.(TypedResolver[M]); err {
			named[nmd.TypedRef().Name()] = nmd
			continue
		}

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

	for name, runner := range named {
		snk.resolvers[name] = runner

		inpts, err := DependancyInputs(name, snk.Resolve, enums...)
		if err != nil {
			return nil, err
		}

		err = opts.Implementation.Decorate(ctx, runner.TypedRef(), snk, inpts)
		if err != nil {
			return nil, err
		}

	}

	err := opts.Implementation.OnSnakeInit(ctx, snk)
	if err != nil {
		return nil, err
	}

	return snk, nil

}
