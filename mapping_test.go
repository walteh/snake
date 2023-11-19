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

func (m MockIsRunnable) HandleResponse(x []reflect.Value) ([]*reflect.Value, error) {
	return []*reflect.Value{&x[0]}, nil
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
		{"key1", []string{"key1", "int"}},
		{"key2", []string{"key2", "int", "uint64", "string"}},
		{"key3", []string{"key3", "int", "uint64"}},
		{"key4", []string{"key4", "int", "uint64", "string", "snake.MockHasRunArgs"}},
	}

	for _, tt := range tableTests {
		t.Run(tt.str, func(t *testing.T) {
			got, err := findBrothers(tt.str, func(s string) HasRunArgs {
				if r, ok := fmap[s]; ok {
					return r
				}
				return nil
			})
			require.NoError(t, err)
			require.NotNil(t, got)
			tt.expectMap = append(tt.expectMap, []string{}...)
			assert.ElementsMatch(t, tt.expectMap, got)
		})
	}
}

func TestFindArguments(t *testing.T) {

	type args struct {
		fmap   map[string]IsRunnable
		target string
	}

	tableTests := []struct {
		name string
		want []any
		args args
	}{
		{
			name: "test1",
			args: args{
				target: "key1",
				fmap: map[string]IsRunnable{
					"uint32": MockIsRunnable{
						fn: reflect.ValueOf(func() (uint32, error) {
							return 2, nil
						}),
					},
					"uint16": MockIsRunnable{
						fn: reflect.ValueOf(func(a uint32) (uint16, error) {
							return uint16(a + 1), nil
						}),
					},
					"int": MockIsRunnable{
						fn: reflect.ValueOf(func(a uint32, b uint16) (int, error) {
							return int(a + uint32(b)), nil
						}),
					},
					"key1": MockIsRunnable{
						fn: reflect.ValueOf(func(a uint32, b uint16, c int) error {
							return nil
						}),
					},
				},
			},
			want: []any{
				(uint32(2)),
				(uint16(3)),
				(int(5)),
				nil, // key1 is a function that returns an error, so we expect nil
			},
		},
	}

	for _, tt := range tableTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findArguments(tt.args.target, func(s string) IsRunnable {
				if r, ok := tt.args.fmap[s]; ok {
					return r
				}
				return nil
			})
			require.NoError(t, err)
			require.NotNil(t, got)
			gotValues := make([]any, len(got))
			for i, v := range got {
				gotValues[i] = v.Interface()
			}
			assert.ElementsMatch(t, tt.want, gotValues)
		})
	}
}
