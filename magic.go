package snake

import (
	"reflect"
)

type inlineResolver[M any] struct {
	self        Runner
	trickref    M
	middlewares []Middleware
	name        string
	description string
}

func (me *inlineResolver[M]) RunFunc() reflect.Value {
	return me.self.RunFunc()
}

func (me *inlineResolver[M]) Ref() Method {
	return me.self.Ref()
}

func (me *inlineResolver[M]) TypedRef() M {
	return me.trickref
}

func (me *inlineResolver[M]) IsResolver() {
}

func (me *inlineResolver[M]) WithMiddleware(mw ...Middleware) TypedResolver[M] {
	me.middlewares = append(me.middlewares, mw...)
	return me
}

func (me *inlineResolver[M]) WithName(name string) TypedResolver[M] {
	me.name = name
	return me
}

func (me *inlineResolver[M]) WithDescription(description string) TypedResolver[M] {
	me.description = description
	return me
}

func (me *inlineResolver[M]) WithRunner(m func() Runner) TypedResolver[M] {
	me.self = m()
	return me
}

func (me *inlineResolver[M]) Name() string {
	return me.name
}

func (me *inlineResolver[M]) Description() string {
	return me.description
}

func NewInlineResolver[M any](typed M, nmd Runner) TypedResolver[M] {
	return &inlineResolver[M]{
		self:        nmd,
		trickref:    typed,
		middlewares: []Middleware{},
	}
}
