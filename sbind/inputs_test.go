package sbind_test

import (
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

type mockInput struct {
	name   string
	shared bool
	parent string
	ptr    any
}

func NewMockInputFromInput(i sbind.Input) *mockInput {
	return &mockInput{
		name:   i.Name(),
		shared: i.Shared(),
		parent: i.Parent(),
		ptr:    i.Ptr(),
	}
}

func TestDependancyInputs(t *testing.T) {

	r1 := &ExampleArgumentResolver{ABC: "abc"}
	vr1 := sbind.MustGetRunMethod(r1)

	r2 := &ExampleÇommand{DEF: "oops"}
	vr2 := sbind.MustGetRunMethod(r2)

	r1d := &DuplicateArgumentResolver{ABC: true}
	vr1d := sbind.MustGetRunMethod(r1d)

	r2d := &DuplicateCommand{}
	vr2d := sbind.MustGetRunMethod(r2d)

	r3 := sbind.NewEnumOptionWithResolver("best-enum-ever", nil, MockEnumA, MockEnumB, MockEnumC)
	vr3 := sbind.MustGetRunMethod(r3)

	m := func(str string) sbind.Resolver {
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
		ptr:    &r1.ABC,
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
		name:   "best-enum-ever",
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

			// require.Equal(t, len(tt.expectedInputs), len(inputs))

			inpts := make([]*mockInput, len(inputs))
			for i, v := range inputs {
				inpts[i] = NewMockInputFromInput(v)
			}

			assert.ElementsMatch(t, tt.expectedInputs, inpts)

			// for _, exp := range tt.expectedInputs {
			// 	var v sbind.Input
			// 	for _, c := range inputs {
			// 		if exp.name == c.Name() && exp.parent == c.Parent() {
			// 			v = c
			// 			break
			// 		}
			// 	}

			// 	require.NotNil(t, v)

			// 	assert.Equal(t, exp.name, v.Name())
			// 	assert.Equal(t, exp.shared, v.Shared())
			// 	assert.Equal(t, exp.parent, v.Parent())
			// 	assert.Equal(t, reflect.ValueOf(exp.ptr).Pointer(), reflect.ValueOf(v.Ptr()).Pointer(), "expected %v, got %v for %s - %s", exp.ptr, v.Ptr(), exp.name, exp.parent)
			// }
		})
	}
}