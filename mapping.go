package snake

import "reflect"

type HasRunArgs interface{ RunArgs() []reflect.Type }

type IsRunnable interface {
	HasRunArgs
	Run() reflect.Value
	HandleResponse([]reflect.Value) (reflect.Value, error)
}

func findBrothers[G HasRunArgs](str string, fmap map[string]G) []string {
	raw := findBrothersRaw(str, fmap, nil)
	resp := make([]string, 0, len(raw))
	for k := range raw {
		if k == str || k == "" {
			continue
		}
		resp = append(resp, k)
	}
	return resp
}

func findBrothersRaw[G HasRunArgs](str string, fmap map[string]G, rmap map[string]bool) map[string]bool {
	if rmap == nil {
		rmap = make(map[string]bool)
	}

	if _, ok := fmap[str]; !ok {
		panic("missing mapping for " + str)
	}

	if rmap[str] {
		return rmap
	}

	rmap[str] = true

	for _, f := range fmap[str].RunArgs() {
		rmap = findBrothersRaw(f.String(), fmap, rmap)
	}

	return rmap
}

func findArguments[G IsRunnable](str string, fmap map[string]G) []reflect.Value {
	raw := findArgumentsRaw(str, fmap, nil)
	resp := make([]reflect.Value, 0, len(raw))
	for _, v := range raw {
		resp = append(resp, v)
	}
	return resp
}

func findArgumentsRaw[G IsRunnable](str string, fmap map[string]G, wrk map[string]reflect.Value) map[string]reflect.Value {

	if _, ok := fmap[str]; !ok {
		panic("missing mapping for " + str)
	}

	if wrk == nil {
		wrk = make(map[string]reflect.Value)
	}

	if _, ok := wrk[str]; ok {
		return wrk
	}

	curr := fmap[str]

	tmp := make([]reflect.Value, len(curr.RunArgs()))
	for _, f := range curr.RunArgs() {
		name := f.String()
		wrk = findArgumentsRaw(name, fmap, wrk)
		tmp = append(tmp, wrk[name])
	}

	resp := curr.Run().Call(tmp)
	var err error
	wrk[str], err = curr.HandleResponse(resp)
	if err != nil {
		panic(err)
	}

	return wrk
}
