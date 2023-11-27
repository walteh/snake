package sbind

import (
	"reflect"

	"github.com/go-faster/errors"
	"github.com/spf13/pflag"
)

// type HasRunArgs interface{ RunArgs() []reflect.Type }

type Method interface {
	// HasRunArgs
	// Run() reflect.Value
	// Names() []string
	// HandleResponse([]reflect.Value) ([]*reflect.Value, error)
}

type Flagged interface {
	Flags() *pflag.FlagSet
}

type FMap[G any] func(string) G

func DependanciesOf[G Method](str string, m FMap[G]) ([]string, error) {
	if ok := m(str); !reflect.ValueOf(ok).IsValid() || reflect.ValueOf(ok).IsNil() {
		return nil, errors.Errorf("missing resolver for %q", str)
	}

	mapa, err := FindBrothers(str, func(s string) Method {
		return m(s)
	})
	if err != nil {
		return nil, err
	}

	return mapa, nil
}

func EndOfChain() reflect.Value {
	return reflect.ValueOf("end_of_chain")
}

func EndOfChainPtr() *reflect.Value {
	v := EndOfChain()
	return &v
}

func FindBrothers(str string, me FMap[Method]) ([]string, error) {
	raw, err := findBrothersRaw(str, me, nil)
	if err != nil {
		return nil, err
	}
	resp := make([]string, 0)
	for k := range raw {
		resp = append(resp, k)
	}
	return resp, nil
}

func findBrothersRaw(str string, fmap FMap[Method], rmap map[string]bool) (map[string]bool, error) {
	var err error
	if rmap == nil {
		rmap = make(map[string]bool)
	}

	var curr reflect.Value

	if meth := fmap(str); meth == nil {
		return nil, errors.Errorf("missing resolver for %q", str)
	} else {
		curr, err = GetRunMethod(meth)
		if err != nil {
			return nil, err
		}
	}

	if rmap[str] {
		return rmap, nil
	}

	rmap[str] = true

	for _, f := range ListOfArgs(curr.Type()) {
		rmap, err = findBrothersRaw(f.String(), fmap, rmap)
		if err != nil {
			return nil, err
		}
	}

	return rmap, nil
}

func FindArguments(str string, fmap FMap[Method]) ([]reflect.Value, error) {
	raw, err := findArgumentsRaw(str, fmap, nil)
	if err != nil {
		return nil, err
	}
	resp := make([]reflect.Value, 0)
	for _, v := range raw.bindings {
		resp = append(resp, *v)
	}
	return resp, nil
}

func valueToMethod(v reflect.Value) Method {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.Interface().(Method)
}

func RunResolvingArguments(str string, fmap FMap[Method]) error {

	_, err := findArgumentsRaw(str, fmap, nil)
	if err != nil {
		return err
	}

	return nil

	// if resp, ok := args.bindings[str]; !ok {
	// 	return errors.Errorf("missing resolver for %q", str)
	// } else {

	// 	// var r reflect.Value

	// 	// if resp.Kind() == reflect.Ptr {
	// 	// 	r = resp.Elem()
	// 	// } else {

	// 	if resp.Interface() != nil {
	// 		// r = resp
	// 		return resp.Interface().(error)
	// 	} else {
	// 		return nil
	// 	}

	// 	// }
	// 	// return resp.Interface().(error)
	// 	// if resp.IsZero() {
	// 	// 	return nil
	// 	// }
	// 	// isError := r.Type().Implements(reflect.TypeOf((*error)(nil)).Elem())
	// 	// if !isError {
	// 	// 	return nil
	// 	// } else {
	// 	// 	return r.Interface().(error)
	// 	// }
	// }

}

func reflectTypeString(typ reflect.Type) string {
	return typ.String()
}

func findArgumentsRaw(str string, fmap FMap[Method], wrk *Binder) (*Binder, error) {
	var curr reflect.Value
	var err error
	if meth := fmap(str); meth == nil {
		return nil, errors.Errorf("missing resolver for %q", str)
	} else {
		curr, err = GetRunMethod(meth)
		if err != nil {
			return nil, err
		}
	}

	if wrk == nil {
		wrk = NewBinder()
	}

	if _, ok := wrk.bindings[str]; ok {
		return wrk, nil
	}

	tmp := make([]reflect.Value, 0)
	for _, f := range ListOfArgs(curr.Type()) {
		name := reflectTypeString(f)
		wrk, err = findArgumentsRaw(name, fmap, wrk)
		if err != nil {
			return nil, err
		}
		tmp = append(tmp, *wrk.bindings[name])
	}

	// var methd reflect.Value

	// if prov, ok := curr.(MethodProvider); ok {
	// 	methd = prov.Method()
	// } else {
	// 	methd = GetRunMethod(curr)
	// }

	out := curr.Call(tmp)
	// out, err := curr.HandleResponse(resp)
	// if err != nil {
	// 	return nil, err
	// }

	if len(out) == 1 {
		// only commands can have one response value, which is always an error
		// so here we know we can name it str
		// otherwise we would be naming it "error"
		wrk.bindings[str] = &out[0]
		if out[0].Interface() != nil {
			return wrk, out[0].Interface().(error)
		}
	} else {
		for _, v := range out {
			in := v
			strd := v.Type().String()
			if strd != "error" {
				wrk.bindings[strd] = &in
			}
		}
	}

	return wrk, nil

}
