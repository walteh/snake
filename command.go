package snake

import (
	"reflect"

	"github.com/go-faster/errors"
)

var (
	ErrMissingResolver = errors.New("missing resolver")
)

func setBindingWithLock[T any](con *Ctx, val T) func() {
	con.runlock.Lock()
	ptr := reflect.ValueOf(val)
	typ := reflect.TypeOf((*T)(nil)).Elem()
	con.bindings[typ.String()] = &ptr
	return func() {
		delete(con.bindings, typ.String())
		con.runlock.Unlock()
	}
}

var (
	ErrInvalidMethodSignature = errors.New("invalid method signatured")
)

var end_of_chain = reflect.ValueOf("end_of_chain")
var end_of_chain_ptr = &end_of_chain

// var prefix_command = "snake_command_"
// var prefix_argument = "snake_argument_"
