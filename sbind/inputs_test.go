package sbind_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/walteh/snake/sbind"
)

type ExampleArgumentResolver struct {
	ABC string
}

func (me *ExampleArgumentResolver) Run() (string, error) {
	return me.ABC, nil
}

type DuplicateArgumentResolver struct {
	ABC bool
}

func (me *DuplicateArgumentResolver) Run() (bool, error) {
	return me.ABC, nil
}

type EnumArgumentResolver struct {
	GHI MockEnum
}

func (me *EnumArgumentResolver) Run() (MockEnum, error) {
	return me.GHI, nil
}

type ExampleÇommand struct {
	DEF string
}

func (me *ExampleÇommand) Run(abc string, en MockEnum) error {
	return nil
}

type DuplicateCommand struct {
}

func (me *DuplicateCommand) Run(abc string, abc2 bool) error {
	return nil
}

type MockEnum string

const (
	MockEnumA MockEnum = "a"
	MockEnumB MockEnum = "b"
	MockEnumC MockEnum = "c"
)

func TestDependancyInputs(t *testing.T) {

	type mockInput struct {
		name   string
		shared bool
		parent string
		ptr    any
	}

	r1 := &ExampleArgumentResolver{
		ABC: "oops",
	}

	r2 := &ExampleÇommand{
		DEF: "oops",
	}

	r1d := &DuplicateArgumentResolver{
		ABC: true,
	}

	r2d := &DuplicateCommand{}

	r3 := &EnumArgumentResolver{
		GHI: MockEnumC,
	}

	m := func(str string) sbind.Method {
		switch str {
		case "bool":
			return r1d
		case "string":
			return r1
		case "sbind_test.MockEnum":
			return r3
		case "command1":
			return r2
		case "command2":
			return r2d
		}
		return nil
	}

	expectedR1 := &mockInput{
		name:   "abc",
		shared: true,
		parent: sbind.MethodName(r1),
		ptr:    &r1.ABC,
	}

	expectedR2 := &mockInput{
		name:   "def",
		shared: false,
		parent: sbind.MethodName(r2),
		ptr:    &r2.DEF,
	}

	expectedR1d := &mockInput{
		name:   "abc",
		shared: true,
		parent: sbind.MethodName(r1d),
		ptr:    &r1d.ABC,
	}

	expectedEnum := &mockInput{
		name:   "ghi",
		shared: true,
		parent: sbind.MethodName(r3),
		ptr:    &r3.GHI,
	}

	tests := []struct {
		name           string
		str            string
		expectedInputs []*mockInput
		wantErr        bool
	}{
		{
			name: "example string",
			str:  "string",
			expectedInputs: []*mockInput{
				expectedR1,
			},
			wantErr: false,
		},
		{
			name: "example command",
			str:  "command1",
			expectedInputs: []*mockInput{
				expectedR2,
				expectedR1,
				expectedEnum,
			},
			wantErr: false,
		},
		{
			name: "example bool",
			str:  "bool",
			expectedInputs: []*mockInput{
				expectedR1d,
			},
			wantErr: false,
		},
		{
			name:           "example command",
			str:            "command2",
			expectedInputs: nil,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputs, err := sbind.DependancyInputs(tt.str, m)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, len(tt.expectedInputs), len(inputs))

			for _, exp := range tt.expectedInputs {
				var v sbind.Input
				for _, c := range inputs {

					if exp.name == c.Name() && exp.parent == c.Parent() {
						v = c
						break
					}
				}

				assert.NotNil(t, v)

				assert.Equal(t, exp.name, v.Name())
				assert.Equal(t, exp.shared, v.Shared())
				assert.Equal(t, exp.parent, v.Parent())
				assert.Equal(t, reflect.ValueOf(exp.ptr).Pointer(), reflect.ValueOf(v.Ptr()).Pointer())
			}
		})
	}
}
