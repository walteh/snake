package sbind

import (
	"context"
	"reflect"
)

type NewSnakeOpts struct {
	Resolvers                  []Method
	NamedResolvers             map[string]Method
	GlobalContextResolverFlags bool
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

// func flagResolver[F any](snk Snake, name string) ([]F, error) {

// 	if flgs, err := DependanciesOf(name, snk.Resolve); err != nil {
// 		return nil, err
// 	} else {
// 		procd := make(map[any]*F, 0)

// 		for _, f := range flgs {
// 			if z := snk.Resolve(f); reflect.ValueOf(z).IsNil() {
// 				return nil, errors.Errorf("missing resolver for %q", f)
// 			} else {
// 				if z, ok := z.(F); ok && procd[z] == nil {
// 					procd[z] = &z
// 				}
// 			}
// 		}
// 		resp := make([]F, 0)
// 		for _, v := range procd {
// 			resp = append(resp, *v)
// 		}

// 		return resp, nil
// 		// fs.VisitAll(func(f *pflag.Flag) {
// 		// 	// if globalFlags != nil && globalFlags.Lookup(f.Name) != nil {
// 		// 	// 	return
// 		// 	// }
// 		// 	cmd.Flags().AddFlag(f)
// 		// })
// 	}

// }

// func ForEachResolvedMethod(snk Snake, fn func(Method)) {
// 	for _, v := range snk.ResolverNames() {
// 		fn(snk.Resolve(v))
// 	}
// }

func NewSnake[M Method](opts *NewSnakeOpts, impl SnakeImplementation[M]) (Snake, error) {

	snk := &defaultSnake{
		bindings:  NewBinder(),
		resolvers: make(map[string]Method),
	}

	// we always want context to get resolved first
	opts.NamedResolvers["root"] = NewNoopAsker[context.Context]()

	for _, v := range opts.Resolvers {

		retrn := ReturnArgs(v)

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
			inpts, err := InputsFor(cmd)
			if err != nil {
				return nil, err
			}

			err = impl.Decorate(cmd, snk, inpts)
			if err != nil {
				return nil, err
			}

			// err = impl.ProcessInputs(cmd, snk)
			// if err != nil {
			// 	return nil, err
			// }

			// err = sbind.NewCommandStrategy().ValidateResponseTypes(sbind.ReturnArgs(exer))
			// if err != nil {
			// 	return nil, err
			// }

			// root.AddCommand(cmd)

			continue
		}

	}

	return snk, nil

}
