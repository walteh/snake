package sbind

import (
	"context"
)

// type Provider[M any] interface {
// 	Provide() M
// }

type NewSnakeOpts struct {
	// Commands                   []Method
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

	// if reflect.TypeOf(res).Kind() == reflect.Ptr {
	// 	if reflect.ValueOf(res).Elem().Kind().String() == "namedMethod" {
	// 		res = reflect.ValueOf(res).Elem().Interface().(MethodWrapper).Method()
	// 	}
	// }

	// for {
	// 	if resd, ok := res.(MethodWrapper); ok {
	// 		maybe := resd.Method()
	// 		if maybed, ok := maybe.(Method); ok {
	// 			res = maybed
	// 		} else {
	// 			break
	// 		}
	// 	} else {
	// 		break
	// 	}
	// }

	// return res
}

func (me *defaultSnake) Binder() *Binder {
	return me.bindings
}

func NewSnake(opts *NewSnakeOpts) (Snake, error) {

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

	return snk, nil

}

type fakeMethodNeedingCongtext struct {
	name string
}

func (me *fakeMethodNeedingCongtext) Run(context.Context) error {
	return nil
}

func (me *fakeMethodNeedingCongtext) Names() []string {
	return []string{me.name}
}

func NewFakeMethodNeedingContext(name string) Method {
	return &fakeMethodNeedingCongtext{name: name}
}
