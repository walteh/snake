package snake_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/walteh/snake"
)

func NewMockIsRunnable(fn any) snake.Resolver {
	return &MockIsRunnable{
		fn: reflect.ValueOf(fn),
	}
}

func (m *MockIsRunnable) Ref() snake.Method {
	return m
}

type MockIsRunnable struct {
	fn reflect.Value
}

func (m *MockIsRunnable) IsResolver() {}

func (m *MockIsRunnable) RunFunc() reflect.Value {
	return m.fn
}

func (m MockIsRunnable) Names() []string {
	return []string{m.fn.Type().String()}
}

func (m MockIsRunnable) HandleResponse(x []reflect.Value) ([]*reflect.Value, error) {
	return []*reflect.Value{&x[0]}, nil
}

func TestFindBrothers(t *testing.T) {
	fmap := map[string]snake.Resolver{
		"int":                       NewMockIsRunnable(func() {}),
		"uint64":                    NewMockIsRunnable(func(int) {}),
		"string":                    NewMockIsRunnable(func(uint64) {}),
		"snake_test.MockIsRunnable": NewMockIsRunnable(func(int, uint64, string) {}),
		"key1":                      NewMockIsRunnable(func(int) {}),
		"key2":                      NewMockIsRunnable(func(int, uint64, string) {}),
		"key3":                      NewMockIsRunnable(func(uint64) {}),
		"key4":                      NewMockIsRunnable(func(MockIsRunnable) {}),
	}

	tableTests := []struct {
		str       string
		expectMap []string
	}{
		{"key1", []string{"key1", "int"}},
		{"key2", []string{"key2", "int", "uint64", "string"}},
		{"key3", []string{"key3", "int", "uint64"}},
		{"key4", []string{"key4", "int", "uint64", "string", "snake_test.MockIsRunnable"}},
	}

	for _, tt := range tableTests {
		t.Run(tt.str, func(t *testing.T) {
			got, err := snake.FindBrothers(tt.str, func(s string) snake.Resolver {
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
		fmap   map[string]snake.Resolver
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
				fmap: map[string]snake.Resolver{
					"uint32": &MockIsRunnable{
						fn: reflect.ValueOf(func() (uint32, error) {
							return 2, nil
						}),
					},
					"uint16": &MockIsRunnable{
						fn: reflect.ValueOf(func(a uint32) (uint16, error) {
							return uint16(a + 1), nil
						}),
					},
					"int": &MockIsRunnable{
						fn: reflect.ValueOf(func(a uint32, b uint16) (int, error) {
							return int(a + uint32(b)), nil
						}),
					},
					"key1": &MockIsRunnable{
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
			got, err := snake.FindArguments(tt.args.target, func(s string) snake.Resolver {
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
