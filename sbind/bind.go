package sbind

import (
	"reflect"
	"sync"
)

type Binder struct {
	bindings map[string]*reflect.Value
	runlock  sync.Mutex
}

func (me *Binder) Bound(name string) *reflect.Value {
	me.runlock.Lock()
	defer me.runlock.Unlock()
	return me.bindings[name]
}

func (me *Binder) Bind(name string, val *reflect.Value) {
	me.runlock.Lock()
	defer me.runlock.Unlock()
	me.bindings[name] = val
}

func NewBinder() *Binder {
	return &Binder{
		bindings: make(map[string]*reflect.Value),
	}
}

func SetBindingWithLock[T any](con *Binder, val T) func() {
	con.runlock.Lock()
	defer con.runlock.Unlock()
	ptr := reflect.ValueOf(val)
	typ := reflect.TypeOf((*T)(nil)).Elem()
	con.bindings[typ.String()] = &ptr
	return func() {
		con.runlock.Lock()
		delete(con.bindings, typ.String())
		con.runlock.Unlock()
	}
}
