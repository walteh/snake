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
	MockHasRunArgs
}

func (m MockIsRunnable) Run() reflect.Value {
	return reflect.ValueOf("result")
}

func (m MockIsRunnable) HandleResponse([]reflect.Value) (reflect.Value, error) {
	return reflect.ValueOf("handled"), nil
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
		"snake.HasRunArgs": MockHasRunArgs{
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
			args: []any{},
		},
	}

	tableTests := []struct {
		str       string
		expectMap []string
	}{
		{"key1", []string{"int"}},
		{"key2", []string{"int", "uint64", "string"}},
		{"key3", []string{"int", "uint64"}},
		{"key4", []string{"int", "uint64", "string", "key2"}},
	}

	for _, tt := range tableTests {
		t.Run(tt.str, func(t *testing.T) {
			got := findBrothers(tt.str, fmap)
			assert.ElementsMatch(t, tt.expectMap, got)
		})
	}
}

func TestFindArguments(t *testing.T) {
	fmap := map[string]IsRunnable{
		"key1": MockIsRunnable{},
		"key2": MockIsRunnable{},
	}

	tableTests := []struct {
		str       string
		expectMap map[string]reflect.Value
	}{
		{"key1", map[string]reflect.Value{"key1": reflect.ValueOf("handled")}},
		{"key2", map[string]reflect.Value{"key1": reflect.ValueOf("handled"), "key2": reflect.ValueOf("handled")}},
	}

	for _, tt := range tableTests {
		t.Run(tt.str, func(t *testing.T) {
			got := findArguments(tt.str, fmap)
			require.NotNil(t, got)
			assert.Equal(t, tt.expectMap, got)
		})
	}
}
