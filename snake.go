package snake

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-faster/errors"
)

type NewSnakeOpts struct {
	Resolvers            []Resolver
	OverrideEnumResolver EnumResolverFunc
}

type Snake interface {
	ResolverNames() []string
	Resolve(string) Resolver
	Enums() []Enum
	Resolvers() []Resolver
	DependantsOf(string) []string
}

type defaultSnake struct {
	resolvers  map[string]Resolver
	dependants map[string][]string
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
	ResolveEnum(string, []string) (string, error)
	ProvideContextResolver() Resolver
}

func NewSnake[M NamedMethod](ctx context.Context, impl SnakeImplementation[M], res ...Resolver) (Snake, error) {
	return NewSnakeWithOpts(ctx, impl, &NewSnakeOpts{
		Resolvers: res,
	})
}

func NewSnakeWithOpts[M NamedMethod](ctx context.Context, impl SnakeImplementation[M], opts *NewSnakeOpts) (Snake, error) {
	var err error

	snk := &defaultSnake{
		resolvers:  make(map[string]Resolver),
		dependants: make(map[string][]string),
	}

	enums := make([]Enum, 0)

	named := make(map[string]TypedResolver[M])

	inputResolvers := make([]Resolver, 0)

	if opts.Resolvers != nil {
		inputResolvers = append(inputResolvers, opts.Resolvers...)
	}

	con := impl.ProvideContextResolver()
	if con != nil {
		inputResolvers = append(inputResolvers, con)
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
			resolver := opts.OverrideEnumResolver
			if resolver == nil {
				resolver = impl.ResolveEnum
			}
			err := mp.ApplyResolver(resolver)
			if err != nil {
				return nil, err
			}
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

	for name := range snk.resolvers {

		deps, err := DependanciesOf(name, snk.Resolve)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to find dependancies of %q", name)
		}

		for _, dep := range deps {
			methd := MethodName(snk.Resolve(dep))

			if _, ok := snk.dependants[methd]; !ok {
				snk.dependants[methd] = make([]string, 0)
			}

			fmt.Println("adding dependant", name, "to", methd)
			snk.dependants[methd] = append(snk.dependants[methd], name)
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

func (me *defaultSnake) Enums() []Enum {
	enums := make([]Enum, 0)
	for _, name := range me.resolvers {
		if mp, ok := name.(Enum); ok {
			enums = append(enums, mp)
		}
	}
	return enums
}

func (me *defaultSnake) Resolvers() []Resolver {
	abc := make([]Resolver, 0)
	for _, name := range me.resolvers {
		abc = append(abc, name)
	}
	return abc
}

func (me *defaultSnake) DependantsOf(name string) []string {
	return me.dependants[name]
}
