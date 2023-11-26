package sbind

import (
	"reflect"

	"github.com/go-faster/errors"
)

type HasRunArgs interface{ RunArgs() []reflect.Type }

type IsRunnable interface {
	HasRunArgs
	Run() reflect.Value
	HandleResponse([]reflect.Value) ([]*reflect.Value, error)
}

type Flagged[F any] interface {
	Flags(*F)
}

type FMap[G any] func(string) G

func FlagsFor[F any, G HasRunArgs](str string, m FMap[G]) (*F, error) {
	if ok := m(str); reflect.ValueOf(ok).IsNil() {
		return nil, errors.Errorf("missing resolver for %q", str)
	}

	mapa, err := FindBrothers(str, func(s string) HasRunArgs {
		return m(s)
	})
	if err != nil {
		return nil, err
	}

	flgs := new(F)

	procd := make(map[Flagged[F]]bool, 0)

	for _, f := range mapa {
		if z := m(f); reflect.ValueOf(z).IsNil() {
			return nil, errors.Errorf("missing resolver for %q", f)
		} else {
			if z, ok := any(z).(Flagged[F]); ok && !procd[z] {
				procd[z] = true
				z.Flags(flgs)
			}
		}
	}

	return flgs, nil
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
// 	args, err := findArgumentsRaw(str, func(s string) IsRunnable {
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

func FindBrothers(str string, me FMap[HasRunArgs]) ([]string, error) {
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

func findBrothersRaw(str string, fmap FMap[HasRunArgs], rmap map[string]bool) (map[string]bool, error) {
	var err error
	if rmap == nil {
		rmap = make(map[string]bool)
	}

	var curr HasRunArgs

	if ok := fmap(str); ok == nil {
		return nil, errors.Errorf("missing resolver for %q", str)
	} else {
		curr = ok
	}

	if rmap[str] {
		return rmap, nil
	}

	rmap[str] = true

	for _, f := range curr.RunArgs() {
		rmap, err = findBrothersRaw(f.String(), fmap, rmap)
		if err != nil {
			return nil, err
		}
	}

	return rmap, nil
}

func FindArguments(str string, fmap FMap[IsRunnable]) ([]reflect.Value, error) {
	raw, err := findArgumentsRaw(str, fmap, nil)
	if err != nil {
		return nil, err
	}
	resp := make([]reflect.Value, 0)
	for _, v := range raw {
		resp = append(resp, *v)
	}
	return resp, nil
}

func valueToIsRunnable(v reflect.Value) IsRunnable {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.Interface().(IsRunnable)
}

func RunResolvingArguments(str string, fmap FMap[IsRunnable], bmap map[string]*reflect.Value) error {

	args, err := findArgumentsRaw(str, fmap, bmap)
	if err != nil {
		return err
	}

	if resp, ok := args[str]; !ok {
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

func findArgumentsRaw(str string, fmap FMap[IsRunnable], wrk map[string]*reflect.Value) (map[string]*reflect.Value, error) {
	var curr IsRunnable
	var err error
	if ok := fmap(str); ok == nil {
		return nil, errors.Errorf("missing resolver for %q", str)
	} else {
		curr = ok
	}

	if wrk == nil {
		wrk = make(map[string]*reflect.Value)
	}

	if _, ok := wrk[str]; ok {
		return wrk, nil
	}

	tmp := make([]reflect.Value, 0)
	for _, f := range curr.RunArgs() {
		name := reflectTypeString(f)
		wrk, err = findArgumentsRaw(name, fmap, wrk)
		if err != nil {
			return nil, err
		}
		tmp = append(tmp, *wrk[name])
	}

	resp := curr.Run().Call(tmp)
	out, err := curr.HandleResponse(resp)
	if err != nil {
		return nil, err
	}

	if len(out) == 1 {
		// only commands can have one response value, which is always an error
		// so here we know we can name it str
		// otherwise we would be naming it "error"
		wrk[str] = out[0]
	} else {
		for _, v := range out {
			if v.Type().String() != "error" {
				wrk[v.Type().String()] = v
			}
		}
	}

	return wrk, nil

}
