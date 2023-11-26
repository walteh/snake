package sbind

import (
	"reflect"
)

type NewSnakeOpts struct {
	Commands                   []Method
	Resolvers                  []Method
	GlobalContextResolverFlags bool
}

type Snake interface {
	Resolve(string) Method
	Bound(string) *reflect.Value
}

type defaultSnake struct {
	bindings  *Binder
	resolvers map[string]Method
}

func NewSnake(opts *NewSnakeOpts) (Snake, error) {

	snk := &defaultSnake{
		bindings:  NewBinder(),
		resolvers: make(map[string]Method),
	}

	for _, v := range opts.Resolvers {
		for _, n := range v.Names() {
			snk.resolvers[n] = v
		}

		if opts.GlobalContextResolverFlags && IsContextResolver(v) {
			v.Flags(root.PersistentFlags())
		}
	}

	for _, v := range opts.Commands {
		for _, n := range v.Names() {
			snk.resolvers[n] = v
		}
	}

}
