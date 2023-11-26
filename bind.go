package snake

import (
	"reflect"
)

func setBindingWithLock[T any](con *Snake, val T) func() {
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
