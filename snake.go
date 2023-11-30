package snake

import (
	"context"
	"reflect"
)

type NewSnakeOpts struct {
	Resolvers []Resolver
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

type SnakeImplementation[X any] interface {
	Decorate(context.Context, X, Snake, []Input, []Middleware) error
	ManagedResolvers(context.Context) []Resolver
	OnSnakeInit(context.Context, Snake) error
}

func NewSnake[M NamedMethod](ctx context.Context, impl SnakeImplementation[M], res ...Resolver) (Snake, error) {
	return NewSnakeWithOpts(ctx, impl, &NewSnakeOpts{
		Resolvers: res,
	})
}

func NewSnakeWithOpts[M NamedMethod](ctx context.Context, impl SnakeImplementation[M], opts *NewSnakeOpts) (Snake, error) {
	var err error

	snk := &defaultSnake{
		resolvers: make(map[string]Resolver),
	}

	enums := make([]Enum, 0)

	named := make(map[string]TypedResolver[M])

	inputResolvers := make([]Resolver, 0)

	if opts.Resolvers != nil {
		inputResolvers = append(inputResolvers, opts.Resolvers...)
	}

	inputResolvers = append(inputResolvers, impl.ManagedResolvers(ctx)...)

	for _, runner := range inputResolvers {

		if nmd, ok := runner.(TypedResolver[M]); ok {
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
		if mp, ok := runner.(Enum); ok {
			enums = append(enums, mp)
		}

	}

	for name, runner := range named {
		snk.resolvers[name] = runner

		inpts, err := DependancyInputs(name, snk.Resolve, enums...)
		if err != nil {
			return nil, err
		}

		mw := make([]Middleware, 0)

		if mwd, ok := runner.(MiddlewareProvider); ok {
			mw = append(mw, mwd.Middlewares()...)

			for _, m := range mwd.Middlewares() {

				mwin, err := InputsFor(NewMiddlewareResolver(m), enums...)
				if err != nil {
					return nil, err
				}

				inpts = append(inpts, mwin...)
			}
		}

		err = impl.Decorate(ctx, runner.TypedRef(), snk, inpts, mw)
		if err != nil {
			return nil, err
		}

	}

	err = impl.OnSnakeInit(ctx, snk)
	if err != nil {
		return nil, err
	}

	return snk, nil

}

func buildMiddlewareName(name string, m Middleware) string {
	return name + "_" + reflectTypeString(reflect.TypeOf(m))
}
