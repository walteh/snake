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

// type EnumArgumentResolver struct {
// 	GHI MockEnum
// }

// func (me *EnumArgumentResolver) Run() (MockEnum, error) {
// 	return me.GHI, nil
// }

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

	r1 := &ExampleArgumentResolver{ABC: "abc"}
	vr1, err := sbind.GetRunMethod(r1)
	require.NoError(t, err)

	r2 := &ExampleÇommand{DEF: "oops"}
	vr2, err := sbind.GetRunMethod(r2)
	require.NoError(t, err)

	r1d := &DuplicateArgumentResolver{ABC: true}
	vr1d, err := sbind.GetRunMethod(r1d)
	require.NoError(t, err)

	r2d := &DuplicateCommand{}
	vr2d, err := sbind.GetRunMethod(r2d)
	require.NoError(t, err)

	r3 := sbind.NewEnumOptionWithResolver(nil, MockEnumA, MockEnumB, MockEnumC)
	vr3, err := sbind.GetRunMethod(r3)
	require.NoError(t, err)

	m := func(str string) sbind.ValidatedRunMethod {
		switch str {
		case "bool":
			return vr1d
		case "string":
			return vr1
		case "sbind_test.MockEnum":
			return vr3
		case "command1":
			return vr2
		case "command2":
			return vr2d
		}
		return nil
	}

	expectedR1 := &mockInput{
		name:   "abc",
		shared: true,
		parent: sbind.MethodName(vr1),
		ptr:    &(vr1.TypedRef()).ABC,
	}

	expectedR2 := &mockInput{
		name:   "def",
		shared: false,
		parent: sbind.MethodName(vr2),
		ptr:    &r2.DEF,
	}

	expectedR1d := &mockInput{
		name:   "abc",
		shared: true,
		parent: sbind.MethodName(vr1d),
		ptr:    &r1d.ABC,
	}

	expectedEnum := &mockInput{
		name:   "myenum",
		shared: true,
		parent: sbind.MethodName(vr3),
		ptr:    r3.CurrentPtr(),
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
			inputs, err := sbind.DependancyInputs(tt.str, m, r3)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, len(tt.expectedInputs), len(inputs))

			for _, exp := range tt.expectedInputs {
				var v sbind.Input
				for _, c := range inputs {

					if exp.name == c.Name() && exp.parent == c.Parent() {
						v = c
						break
					}
				}

				require.NotNil(t, v)

				assert.Equal(t, exp.name, v.Name())
				assert.Equal(t, exp.shared, v.Shared())
				assert.Equal(t, exp.parent, v.Parent())
				assert.Equal(t, reflect.ValueOf(exp.ptr).Pointer(), reflect.ValueOf(v.Ptr()).Pointer())
			}
		})
	}
}
