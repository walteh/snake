package sbind

import (
	"reflect"

	"github.com/go-faster/errors"
)

type Method interface {
}

type FMap func(string) ValidatedRunMethod

func DependanciesOf(str string, m FMap) ([]string, error) {
	if ok := m(str); !reflect.ValueOf(ok).IsValid() || reflect.ValueOf(ok).IsNil() {
		return nil, errors.Errorf("missing resolver for %q", str)
	}

	mapa, err := FindBrothers(str, m)
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

func FindBrothers(str string, me FMap) ([]string, error) {
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

func findBrothersRaw(str string, fmap FMap, rmap map[string]bool) (map[string]bool, error) {
	var err error
	if rmap == nil {
		rmap = make(map[string]bool)
	}

	validated := fmap(str)

	if validated == nil {
		return nil, errors.Errorf("missing resolver for %q", str)
	}

	if rmap[str] {
		return rmap, nil
	}

	rmap[str] = true

	for _, f := range ListOfArgs(validated) {
		rmap, err = findBrothersRaw(f.String(), fmap, rmap)
		if err != nil {
			return nil, err
		}
	}

	return rmap, nil
}

func FindArguments(str string, fmap FMap) ([]reflect.Value, error) {
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

func RunResolvingArguments(str string, fmap FMap) error {
	_, err := findArgumentsRaw(str, fmap, nil)
	if err != nil {
		return err
	}

	return nil
}

func reflectTypeString(typ reflect.Type) string {
	return typ.String()
}

func findArgumentsRaw(str string, fmap FMap, wrk *Binder) (*Binder, error) {
	validated := fmap(str)
	var err error
	if validated == nil {
		return nil, errors.Errorf("missing resolver for %q", str)
	}

	if wrk == nil {
		wrk = NewBinder()
	}

	if _, ok := wrk.bindings[str]; ok {
		return wrk, nil
	}

	tmp := make([]reflect.Value, 0)
	for _, f := range ListOfArgs(validated) {
		name := reflectTypeString(f)
		wrk, err = findArgumentsRaw(name, fmap, wrk)
		if err != nil {
			return nil, err
		}
		tmp = append(tmp, *wrk.bindings[name])
	}

	out := CallMethod(validated, tmp)

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
