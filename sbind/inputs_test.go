package sbind

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
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

type ExampleÇommand struct {
	DEF string
}

func (me *ExampleÇommand) Run(abc string) error {
	return nil
}

type DuplicateCommand struct {
}

func (me *DuplicateCommand) Run(abc string, abc2 bool) error {
	return nil
}

func TestDependancyInputs(t *testing.T) {

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

	m := func(str string) Method {
		switch str {
		case "bool":
			return r1d
		case "string":
			return r1
		case "command1":
			return r2
		case "command2":
			return r2d
		}
		return nil
	}

	expectedR1 := &MockInput{
		name:   "abc",
		typ:    reflect.TypeOf(r1.ABC),
		shared: true,
		m:      r1,
		val:    &r1.ABC,
		usage:  "",
	}

	expectedR2 := &MockInput{
		name:   "def",
		typ:    reflect.TypeOf(r2.DEF),
		shared: false,
		m:      r2,
		val:    &r2.DEF,
		usage:  "",
	}

	expectedR1d := &MockInput{
		name:   "abc",
		typ:    reflect.TypeOf(r1d.ABC),
		shared: true,
		m:      r1d,
		val:    &r1d.ABC,
		usage:  "",
	}

	tests := []struct {
		name           string
		str            string
		expectedInputs []Input
		wantErr        bool
	}{
		{
			name: "example string",
			str:  "string",
			expectedInputs: []Input{
				expectedR1,
			},
			wantErr: false,
		},
		{
			name: "example command",
			str:  "command1",
			expectedInputs: []Input{
				expectedR2,
				expectedR1,
			},
			wantErr: false,
		},
		{
			name: "example bool",
			str:  "bool",
			expectedInputs: []Input{
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
			inputs, err := DependancyInputs(tt.str, m)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, len(tt.expectedInputs), len(inputs))

			for _, exp := range tt.expectedInputs {
				var v Input
				for _, c := range tt.expectedInputs {
					if exp.Ptr() == c.Ptr() {
						v = c
						break
					}
				}
				assert.NotNil(t, exp)
				assert.Equal(t, exp.Name(), v.Name())
				assert.Equal(t, exp.Type().Name(), v.Type().Name())
				assert.Equal(t, exp.Shared(), v.Shared())
				assert.Equal(t, exp.M(), v.M())
				assert.Equal(t, exp.Usage(), v.Usage())
				assert.Equal(t, exp.Ptr(), v.Ptr())
			}
		})
	}
}

type MockInput struct {
	name   string
	typ    reflect.Type
	shared bool
	m      Method
	val    any
	usage  string
}

func (me *MockInput) Name() string {
	return me.name
}

func (me *MockInput) Type() reflect.Type {
	return me.typ
}

func (me *MockInput) Shared() bool {
	return me.shared
}

func (me *MockInput) M() Method {
	return me.m
}

func (me *MockInput) Usage() string {
	return me.usage
}

func (me *MockInput) Ptr() any {
	return me.val
}
