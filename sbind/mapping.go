package sbind

import (
	"fmt"
	"reflect"

	"github.com/go-faster/errors"
	"github.com/spf13/pflag"
)

// type HasRunArgs interface{ RunArgs() []reflect.Type }

type Method interface {
	// HasRunArgs
	Run() reflect.Value
	Names() []string
	// HandleResponse([]reflect.Value) ([]*reflect.Value, error)
}

type Flagged interface {
	Flags() *pflag.FlagSet
}

type FMap[G any] func(string) G

func FlagsFor[G Method](str string, m FMap[G]) ([]string, error) {
	if ok := m(str); reflect.ValueOf(ok).IsNil() {
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

// func (me *Snake) Run(str Method) error {
// 	return me.RunString(str.Name())
// }

func EndOfChain() reflect.Value {
	return reflect.ValueOf("end_of_chain")
}

func EndOfChainPtr() *reflect.Value {
	v := EndOfChain()
	return &v
}

// func (me *Snake) RunString(str string) error {
// 	args, err := findArgumentsRaw(str, func(s string) Method {
// 		return me.resolvers[s]
// 	}, nil)
// 	if err != nil {
// 		return err
// 	}

// 	if resp, ok := args[str]; !ok {
// 		return errors.Wrapf(ErrMissingResolver, "missing resolver for %q", str)
// 	} else {
// 		if resp == end_of_chain_ptr {
// 			return nil
// 		} else {
// 			return errors.Errorf("expected end of chain, got %v", resp)
// 		}
// 	}

// }

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

	var curr Method

	if ok := fmap(str); ok == nil {
		return nil, errors.Errorf("missing resolver for %q", str)
	} else {
		curr = ok
	}

	if rmap[str] {
		return rmap, nil
	}

	rmap[str] = true

	for _, f := range RunArgs(curr) {
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

	args, err := findArgumentsRaw(str, fmap, NewBinder())
	if err != nil {
		return err
	}

	if resp, ok := args.bindings[str]; !ok {
		return errors.Errorf("missing resolver for %q", str)
	} else {
		if reflect.DeepEqual(resp, EndOfChainPtr()) {
			return nil
		} else {
			return errors.Errorf("expected end of chain, got %v", resp)
		}
	}

}

func reflectTypeString(typ reflect.Type) string {
	return typ.String()
}

func findArgumentsRaw(str string, fmap FMap[Method], wrk *Binder) (*Binder, error) {
	var curr Method
	var err error
	if ok := fmap(str); ok == nil {
		return nil, errors.Errorf("missing resolver for %q", str)
	} else {
		curr = ok
	}

	if wrk == nil {
		wrk = NewBinder()
	}

	if _, ok := wrk.bindings[str]; ok {
		return wrk, nil
	}

	tmp := make([]reflect.Value, 0)
	for _, f := range RunArgs(curr) {
		name := reflectTypeString(f)
		wrk, err = findArgumentsRaw(name, fmap, wrk)
		if err != nil {
			return nil, err
		}
		tmp = append(tmp, *wrk.bindings[name])
	}

	out := curr.Run().Call(tmp)
	// out, err := curr.HandleResponse(resp)
	// if err != nil {
	// 	return nil, err
	// }

	if len(out) == 1 {
		// only commands can have one response value, which is always an error
		// so here we know we can name it str
		// otherwise we would be naming it "error"
		wrk.bindings[str] = &out[0]
	} else {
		for _, v := range out {
			in := v
			fmt.Println(v.Type().String())
			strd := v.Type().String()
			if strd != "error" {
				wrk.bindings[strd] = &in
			}
		}
	}

	return wrk, nil

}
