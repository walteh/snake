package snake

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestBind(t *testing.T) {
	ctx := context.Background()
	key := &struct{}{}
	value := "value"

	newCtx := Bind(ctx, key, value)
	b, ok := newCtx.Value(&bindingsKeyT{}).(bindings)

	assert.True(t, ok)
	assert.NotNil(t, b)
}

func TestCallRunMethod(t *testing.T) {
	ctx := context.Background()

	cmd := &cobra.Command{}
	f := reflect.ValueOf(func(ctx context.Context) error {
		return nil
	})

	err := callRunMethod(ctx, cmd, f, f.Type())

	assert.Nil(t, err)
}

func TestNewRootCommandNoRun(t *testing.T) {

	snakeableMock := &MockSnakeableNoRun{
		ParseArgumentsFunc: func(ctx context.Context, cmd *cobra.Command, args []string) error {
			return nil
		},
		BuildCommandFunc: func(ctx context.Context) *cobra.Command {
			return &cobra.Command{}
		},
	}

	ctx := context.Background()

	cmd := NewRootCommand(ctx, snakeableMock)

	assert.NotNil(t, cmd)
}

func TestNewCommandNoRun(t *testing.T) {
	// Set up your Snakeable mock
	snakeableMock := &MockSnakeableNoRun{
		ParseArgumentsFunc: func(ctx context.Context, cmd *cobra.Command, args []string) error {
			return nil
		},
		BuildCommandFunc: func(ctx context.Context) *cobra.Command {
			return &cobra.Command{}
		},
	}

	rootCmd := &cobra.Command{}
	ctx := context.Background()

	err := NewCommand(ctx, rootCmd, snakeableMock)

	assert.ErrorIs(t, err, ErrMissingRun)
}

func TestNewCommandValid(t *testing.T) {
	// Set up your Snakeable mock
	snakeableMock := &MockSnakeableWithZeroInput{
		MockSnakeableNoRun: MockSnakeableNoRun{
			ParseArgumentsFunc: func(ctx context.Context, cmd *cobra.Command, args []string) error {
				return nil
			},
			BuildCommandFunc: func(ctx context.Context) *cobra.Command {
				return &cobra.Command{}
			},
		},
		RunFunc: func() error {
			return nil
		},
	}

	rootCmd := &cobra.Command{}
	ctx := context.Background()

	err := NewCommand(ctx, rootCmd, snakeableMock)

	assert.Nil(t, err)
}

type MockSnakeableCase interface {
	Snakeable
	ExpectedNewCommandError() error
	ExpectedRunCommandError() error

	Bindings() []any
}

// //////////////////////////////////////////////////////////////

var _ MockSnakeableCase = (*MockSnakeableNoRun)(nil)

type MockSnakeableNoRun struct {
	ParseArgumentsFunc func(ctx context.Context, cmd *cobra.Command, args []string) error
	BuildCommandFunc   func(ctx context.Context) *cobra.Command
}

func NewMockSnakeableNoRun() *MockSnakeableNoRun {
	return &MockSnakeableNoRun{
		ParseArgumentsFunc: func(ctx context.Context, cmd *cobra.Command, args []string) error {
			return nil
		},
		BuildCommandFunc: func(ctx context.Context) *cobra.Command {
			return &cobra.Command{
				Use: "mock",
			}
		},
	}
}

func (m *MockSnakeableNoRun) ParseArguments(ctx context.Context, cmd *cobra.Command, args []string) error {
	return m.ParseArgumentsFunc(ctx, cmd, args)
}

func (m *MockSnakeableNoRun) BuildCommand(ctx context.Context) *cobra.Command {
	return m.BuildCommandFunc(ctx)
}

func (m *MockSnakeableNoRun) ExpectedNewCommandError() error {
	return ErrMissingRun
}

func (m *MockSnakeableNoRun) ExpectedRunCommandError() error {
	return nil
}

func (m *MockSnakeableNoRun) Bindings() []any {
	return []any{}
}

////////////////////////////////////////////////////////////////

var _ MockSnakeableCase = (*MockSnakeableWithZeroInput)(nil)

type MockSnakeableWithZeroInput struct {
	MockSnakeableNoRun
	RunFunc func() error
}

func (m *MockSnakeableWithZeroInput) Run() error {
	return m.RunFunc()
}

func (m *MockSnakeableWithZeroInput) ExpectedNewCommandError() error {
	return nil
}

////////////////////////////////////////////////////////////////

type SnakeableWithZeroInputTwoOutput struct {
	MockSnakeableNoRun
	RunFunc func() (string, error)
}

func (m *SnakeableWithZeroInputTwoOutput) Run() (string, error) {
	return m.RunFunc()
}

func (m *SnakeableWithZeroInputTwoOutput) ExpectedNewCommandError() error {
	return ErrInvalidRun
}

////////////////////////////////////////////////////////////////

type SnakeableWithOneInput struct {
	MockSnakeableNoRun
	RunFunc func(context.Context) error
}

func (m *SnakeableWithOneInput) Run(ctx context.Context) error {
	return m.RunFunc(ctx)
}

func (m *SnakeableWithOneInput) ExpectedNewCommandError() error {
	return nil
}

////////////////////////////////////////////////////////////////

type SnakeableWithTwoInput struct {
	MockSnakeableNoRun
	RunFunc func(context.Context, string) error
}

func (m *SnakeableWithTwoInput) Run(ctx context.Context, s string) error {
	return m.RunFunc(ctx, s)
}

func (m *SnakeableWithTwoInput) ExpectedNewCommandError() error {
	return nil
}

func (m *SnakeableWithTwoInput) Bindings() []any {
	return []any{"hi"}
}

////////////////////////////////////////////////////////////////

type SnakeableWithTwoInputContextSecond struct {
	MockSnakeableNoRun
	RunFunc func(string, context.Context) error
}

func (m *SnakeableWithTwoInputContextSecond) Run(s string, ctx context.Context) error {
	return m.RunFunc(s, ctx)
}

func (m *SnakeableWithTwoInputContextSecond) ExpectedNewCommandError() error {
	return nil
}

func (m *SnakeableWithTwoInputContextSecond) Bindings() []any {
	return []any{"hi"}
}

////////////////////////////////////////////////////////////////

type SnakeableWithTwoInputMissingBinding struct {
	MockSnakeableNoRun
	RunFunc func(context.Context, string) error
}

func (m *SnakeableWithTwoInputMissingBinding) Run(ctx context.Context, s string) error {
	return m.RunFunc(ctx, s)
}

func (m *SnakeableWithTwoInputMissingBinding) Bindings() []any {
	return []any{int64(1)}
}

func (m *SnakeableWithTwoInputMissingBinding) ExpectedNewCommandError() error {
	return nil
}

func (m *SnakeableWithTwoInputMissingBinding) ExpectedRunCommandError() error {
	return ErrMissingBinding
}

////////////////////////////////////////////////////////////////

type SnakeableWithThreeInputCobraPointer struct {
	MockSnakeableNoRun
	RunFunc func(context.Context, string, *cobra.Command) error
}

func (m *SnakeableWithThreeInputCobraPointer) Run(ctx context.Context, s string, cmd *cobra.Command) error {
	return m.RunFunc(ctx, s, cmd)
}

func (m *SnakeableWithThreeInputCobraPointer) ExpectedNewCommandError() error {
	return nil
}

func (m *SnakeableWithThreeInputCobraPointer) Bindings() []any {
	return []any{"hi"}
}

////////////////////////////////////////////////////////////////

type SnakeableWithThreeInputCobraNonPointer struct {
	MockSnakeableNoRun
	RunFunc func(context.Context, string, cobra.Command) error
}

func (m *SnakeableWithThreeInputCobraNonPointer) Run(ctx context.Context, s string, cmd cobra.Command) error {
	return m.RunFunc(ctx, s, cmd)
}

func (m *SnakeableWithThreeInputCobraNonPointer) ExpectedNewCommandError() error {
	return nil
}

func (m *SnakeableWithThreeInputCobraNonPointer) Bindings() []any {
	return []any{"hi"}
}

func TestGetRunMethodNoBindings(t *testing.T) {
	tests := []MockSnakeableCase{
		NewMockSnakeableNoRun(),
		&MockSnakeableWithZeroInput{
			MockSnakeableNoRun: *NewMockSnakeableNoRun(),
			RunFunc:            func() error { return nil },
		},
		&SnakeableWithOneInput{
			MockSnakeableNoRun: *NewMockSnakeableNoRun(),
			RunFunc:            func(context.Context) error { return nil },
		},
		&SnakeableWithTwoInput{
			MockSnakeableNoRun: *NewMockSnakeableNoRun(),
			RunFunc:            func(context.Context, string) error { return nil },
		},
		&SnakeableWithTwoInputContextSecond{
			MockSnakeableNoRun: *NewMockSnakeableNoRun(),
			RunFunc:            func(string, context.Context) error { return nil },
		},
		&SnakeableWithZeroInputTwoOutput{
			MockSnakeableNoRun: *NewMockSnakeableNoRun(),
			RunFunc:            func() (string, error) { return "", nil },
		},
		&SnakeableWithTwoInputMissingBinding{
			MockSnakeableNoRun: *NewMockSnakeableNoRun(),
			RunFunc:            func(context.Context, string) error { return nil },
		},
		&SnakeableWithThreeInputCobraPointer{
			MockSnakeableNoRun: *NewMockSnakeableNoRun(),
			RunFunc:            func(context.Context, string, *cobra.Command) error { return nil },
		},
		&SnakeableWithThreeInputCobraNonPointer{
			MockSnakeableNoRun: *NewMockSnakeableNoRun(),
			RunFunc:            func(context.Context, string, cobra.Command) error { return nil },
		},
	}

	for _, tt := range tests {
		t.Run(reflect.ValueOf(tt).String(), func(t *testing.T) {

			ctx := context.Background()

			cmd := cobra.Command{}

			err := NewCommand(ctx, &cmd, tt)
			assert.ErrorIs(t, err, tt.ExpectedNewCommandError())

			if err != nil {
				return
			}

			for _, b := range tt.Bindings() {
				ctx = Bind(ctx, reflect.ValueOf(b).Interface(), b)
			}

			os.Args = []string{"x", "mock"}

			err = cmd.ExecuteContext(ctx)
			assert.ErrorIs(t, err, tt.ExpectedRunCommandError())
		})
	}
}
