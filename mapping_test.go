package snake

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockHasRunArgs struct {
	args []any
}

func (m MockHasRunArgs) RunArgs() []reflect.Type {
	wrk := make([]reflect.Type, len(m.args))
	for i, v := range m.args {
		wrk[i] = reflect.TypeOf(v)
	}
	return wrk
}

type MockIsRunnable struct {
	fn reflect.Value
}

func (m MockIsRunnable) RunArgs() []reflect.Type {
	return listOfArgs(m.fn.Type())
}

func (m MockIsRunnable) Run() reflect.Value {
	return m.fn
}

func (m MockIsRunnable) HandleResponse(x []reflect.Value) (reflect.Value, error) {
	return x[0], nil
}

func TestFindBrothers(t *testing.T) {
	fmap := map[string]HasRunArgs{
		"int": MockHasRunArgs{
			args: []any{},
		},
		"uint64": MockHasRunArgs{
			args: []any{1},
		},
		"string": MockHasRunArgs{
			args: []any{uint64(1)},
		},
		"snake.MockHasRunArgs": MockHasRunArgs{
			args: []any{1, uint64(1), "1"},
		},
		"key1": MockHasRunArgs{
			args: []any{1},
		},
		"key2": MockHasRunArgs{
			args: []any{1, uint64(1), "1"},
		},
		"key3": MockHasRunArgs{
			args: []any{uint64(1)},
		},
		"key4": MockHasRunArgs{
			args: []any{MockHasRunArgs{}},
		},
	}

	tableTests := []struct {
		str       string
		expectMap []string
	}{
		{"key1", []string{"int"}},
		{"key2", []string{"int", "uint64", "string"}},
		{"key3", []string{"int", "uint64"}},
		{"key4", []string{"int", "uint64", "string", "snake.MockHasRunArgs"}},
	}

	for _, tt := range tableTests {
		t.Run(tt.str, func(t *testing.T) {
			got := findBrothers(tt.str, fmap)
			require.NotNil(t, got)
			assert.ElementsMatch(t, tt.expectMap, got)
		})
	}
}

func TestFindArguments(t *testing.T) {
	fmap := map[string]IsRunnable{
		"key1": MockIsRunnable{
			fn: reflect.ValueOf(func() uint32 {
				return 2
			}),
		},
		"key2": MockIsRunnable{
			fn: reflect.ValueOf(func(a uint32) uint16 {
				return uint16(a + 1)
			}),
		},
	}

	tableTests := []struct {
		str       string
		expectMap []any
	}{
		{"key2", []any{uint32(2)}},
	}

	for _, tt := range tableTests {
		t.Run(tt.str, func(t *testing.T) {
			got := findArguments(tt.str, fmap)
			gotres := make([]any, len(got))
			for i, v := range got {
				gotres[i] = v.Interface()
			}
			require.NotNil(t, got)
			assert.ElementsMatch(t, tt.expectMap, gotres)
		})
	}
}
