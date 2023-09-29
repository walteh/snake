package snake

import (
	"context"
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	cmd.SetContext(ctx)

	err := callRunMethod(cmd, f, f.Type())

	assert.Nil(t, err)
}

func TestNewRootCommandNoRun(t *testing.T) {

	snakeableMock := &MockSnakeableNoRun{
		ParseArgumentsFunc: func(ctx context.Context, args []string) (context.Context, error) {
			return ctx, nil
		},
		BuildCommandFunc: func(ctx context.Context) (*cobra.Command, error) {
			return &cobra.Command{}, nil
		},
	}

	ctx, err := NewRootCommand(context.Background(), snakeableMock)
	assert.Nil(t, err)

	cmd := GetRootCommand(ctx)

	assert.NotNil(t, cmd)
}

func TestNewCommandNoRun(t *testing.T) {
	// Set up your Snakeable mock
	snakeableMock := &MockSnakeableNoRun{
		ParseArgumentsFunc: func(ctx context.Context, args []string) (context.Context, error) {
			return ctx, nil
		},
		BuildCommandFunc: func(ctx context.Context) (*cobra.Command, error) {
			return &cobra.Command{}, nil
		},
	}

	ctx, err := NewRootCommand(context.Background(), snakeableMock)
	assert.Nil(t, err)

	_, err = NewCommand(ctx, "hi", snakeableMock)
	assert.ErrorIs(t, err, ErrMissingRun)
}

func TestNewCommandValid(t *testing.T) {
	// Set up your Snakeable mock
	snakeableMock := &MockSnakeableWithZeroInput{
		MockSnakeableNoRun: MockSnakeableNoRun{
			ParseArgumentsFunc: func(ctx context.Context, args []string) (context.Context, error) {
				return ctx, nil
			},
			BuildCommandFunc: func(ctx context.Context) (*cobra.Command, error) {
				return &cobra.Command{}, nil
			},
		},
		RunFunc: func() error {
			return nil
		},
	}

	ctx, err := NewRootCommand(context.Background(), snakeableMock)
	assert.Nil(t, err)

	_, err = NewCommand(ctx, "hi", snakeableMock)

	assert.Nil(t, err)
}

type MockSnakeableCase interface {
	Snakeable
	ExpectedNewCommandError() error
	ExpectedRunCommandError() error

	Bindings() []any
	RootParseArgumentsBindings() []any
}

type customStruct struct {
}

type customInterface interface {
	ID() string
}

func (c *customStruct) ID() string {
	return "id"
}

// //////////////////////////////////////////////////////////////

var _ MockSnakeableCase = (*MockSnakeableNoRun)(nil)

type MockSnakeableNoRun struct {
	ParseArgumentsFunc func(ctx context.Context, args []string) (context.Context, error)
	BuildCommandFunc   func(ctx context.Context) (*cobra.Command, error)
	ResolveBindingFunc func(context.Context, any) (any, error)
}

func NewMockSnakeableNoRun() *MockSnakeableNoRun {
	return &MockSnakeableNoRun{
		ParseArgumentsFunc: func(ctx context.Context, args []string) (context.Context, error) {
			return ctx, nil
		},
		BuildCommandFunc: func(ctx context.Context) (*cobra.Command, error) {
			return &cobra.Command{
				Use: "mock",
			}, nil
		},
		ResolveBindingFunc: nil,
	}
}

func NewMockSnakeableResolveBinding(f func(context.Context, any) (any, error)) *MockSnakeableNoRun {
	snk := NewMockSnakeableNoRun()
	return &MockSnakeableNoRun{
		ParseArgumentsFunc: snk.ParseArgumentsFunc,
		BuildCommandFunc:   snk.BuildCommandFunc,
		ResolveBindingFunc: f,
	}
}

func (m *MockSnakeableNoRun) PreRun(ctx context.Context, args []string) (context.Context, error) {
	return m.ParseArgumentsFunc(ctx, args)
}

func (m *MockSnakeableNoRun) ResolveBinding(ctx context.Context, arg any) (any, error) {
	if m.ResolveBindingFunc == nil {
		return nil, errors.New("no binding resolver")
	}
	return m.ResolveBindingFunc(ctx, arg)
}

func (m *MockSnakeableNoRun) Register(ctx context.Context) (*cobra.Command, error) {
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

func (m *MockSnakeableNoRun) RootParseArgumentsBindings() []any {
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

////////////////////////////////////////////////////////////////

type SnakeableWithCustomStruct struct {
	MockSnakeableNoRun
	RunFunc func(customStruct) error
}

func (m *SnakeableWithCustomStruct) Run(cs customStruct) error {
	return m.RunFunc(cs)
}

func (m *SnakeableWithCustomStruct) ExpectedNewCommandError() error {
	return nil
}

func (m *SnakeableWithCustomStruct) RootParseArgumentsBindings() []any {
	return []any{customStruct{}}
}

////////////////////////////////////////////////////////////////

type SnakeableWithCustomStructPtr struct {
	MockSnakeableNoRun
	RunFunc func(*customStruct) error
}

func (m *SnakeableWithCustomStructPtr) Run(cs *customStruct) error {
	return m.RunFunc(cs)
}

func (m *SnakeableWithCustomStructPtr) ExpectedNewCommandError() error {
	return nil
}

func (m *SnakeableWithCustomStructPtr) RootParseArgumentsBindings() []any {
	return []any{&customStruct{}}
}

////////////////////////////////////////////////////////////////

type SnakeableWithCustomStructPtrInvalidBinding struct {
	MockSnakeableNoRun
	RunFunc func(*customStruct) error
}

func (m *SnakeableWithCustomStructPtrInvalidBinding) Run(cs *customStruct) error {
	return m.RunFunc(cs)
}

func (m *SnakeableWithCustomStructPtrInvalidBinding) ExpectedNewCommandError() error {
	return nil
}

func (m *SnakeableWithCustomStructPtrInvalidBinding) RootParseArgumentsBindings() []any {
	return []any{customStruct{}}
}

func (m *SnakeableWithCustomStructPtrInvalidBinding) ExpectedRunCommandError() error {
	return ErrMissingBinding
}

////////////////////////////////////////////////////////////////

type SnakeableWithBindingFuncNoBindings struct {
	MockSnakeableNoRun
	RunFunc func(*customStruct) error
}

func (m *SnakeableWithBindingFuncNoBindings) Run(cs *customStruct) error {
	return m.RunFunc(cs)
}

func (m *SnakeableWithBindingFuncNoBindings) ExpectedNewCommandError() error {
	return nil
}

func (m *SnakeableWithBindingFuncNoBindings) RootParseArgumentsBindings() []any {
	return []any{}
}

func (m *SnakeableWithBindingFuncNoBindings) ExpectedRunCommandError() error {
	return nil
}

////////////////////////////////////////////////////////////////

type SnakeableWithCustomInterfaceRunFunc struct {
	MockSnakeableNoRun
	RunFunc func(customInterface) error
}

func (m *SnakeableWithCustomInterfaceRunFunc) Run(cs customInterface) error {
	return m.RunFunc(cs)
}

func (m *SnakeableWithCustomInterfaceRunFunc) ExpectedNewCommandError() error {
	return nil
}

func (m *SnakeableWithCustomInterfaceRunFunc) RootParseArgumentsBindings() []any {
	return []any{}
}

func (m *SnakeableWithCustomInterfaceRunFunc) ExpectedRunCommandError() error {
	return nil
}

////////////////////////////////////////////////////////////////

type SnakeableWithCustomInterfaceRunFuncContext struct {
	MockSnakeableNoRun
	RunFunc func(context.Context, customInterface) error
}

func (m *SnakeableWithCustomInterfaceRunFuncContext) Run(ctx context.Context, cs customInterface) error {
	return m.RunFunc(ctx, cs)
}

func (m *SnakeableWithCustomInterfaceRunFuncContext) ExpectedNewCommandError() error {
	return nil
}

func (m *SnakeableWithCustomInterfaceRunFuncContext) RootParseArgumentsBindings() []any {
	return []any{}
}

func (m *SnakeableWithCustomInterfaceRunFuncContext) ExpectedRunCommandError() error {
	return nil
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
		// &SnakeableWithCustomStruct{
		// 	MockSnakeableNoRun: *NewMockSnakeableNoRun(),
		// 	RunFunc:            func(customStruct) error { return nil },
		// },
		// &SnakeableWithCustomStructPtr{
		// 	MockSnakeableNoRun: *NewMockSnakeableNoRun(),
		// 	RunFunc:            func(*customStruct) error { return nil },
		// },
	}

	for _, tt := range tests {
		t.Run(reflect.ValueOf(tt).String(), func(t *testing.T) {

			ctx := context.Background()

			rootcmd := NewMockSnakeableNoRun()

			rootcmd.ParseArgumentsFunc = func(ctx context.Context, args []string) (context.Context, error) {
				for _, b := range tt.RootParseArgumentsBindings() {
					ctx = Bind(ctx, reflect.ValueOf(b).Interface(), b)
				}
				// cmd.SetContext(ctx)
				return ctx, nil
			}

			var err error

			ctx, err = NewRootCommand(ctx, rootcmd)

			ctx, err = NewCommand(ctx, "hello123", tt)
			assert.ErrorIs(t, err, tt.ExpectedNewCommandError())

			if err != nil {
				return
			}

			for _, b := range tt.Bindings() {
				ctx = Bind(ctx, reflect.ValueOf(b).Interface(), b)
			}

			os.Args = []string{"x", "hello123"}

			assm := Assemble(ctx)

			err = assm.ExecuteContext(ctx)

			assert.ErrorIs(t, err, tt.ExpectedRunCommandError())
		})
	}
}

func TestGetRunMethodWithBindingResolverRegistered(t *testing.T) {

	tt := &SnakeableWithBindingFuncNoBindings{
		MockSnakeableNoRun: *NewMockSnakeableNoRun(),
		RunFunc:            func(*customStruct) error { return nil },
	}

	t.Run(reflect.ValueOf(tt).String(), func(t *testing.T) {

		ctx := context.Background()

		rootcmd := NewMockSnakeableNoRun()

		rootcmd.ParseArgumentsFunc = func(ctx context.Context, args []string) (context.Context, error) {
			for _, b := range tt.RootParseArgumentsBindings() {
				ctx = Bind(ctx, reflect.ValueOf(b).Interface(), b)
			}

			// cmd.SetContext(ctx)
			return ctx, nil
		}
		var err error
		ctx, err = NewRootCommand(ctx, rootcmd)
		if err != nil {
			t.Fatal(err)
		}

		ctx = RegisterBindingResolver(ctx, func(context.Context) (*customStruct, error) {
			cms := customStruct{}
			return &cms, nil
		})

		ctx, err = NewCommand(ctx, "hello123", tt)
		assert.ErrorIs(t, err, tt.ExpectedNewCommandError())

		if err != nil {
			return
		}

		for _, b := range tt.Bindings() {
			ctx = Bind(ctx, reflect.ValueOf(b).Interface(), b)
		}

		os.Args = []string{"x", "hello123"}

		assm := Assemble(ctx)

		err = assm.ExecuteContext(ctx)

		assert.ErrorIs(t, err, tt.ExpectedRunCommandError())
	})
}

func TestGetRunMethodWithBindingResolverRegisteredInterfacePtr(t *testing.T) {

	tt := &SnakeableWithCustomInterfaceRunFunc{
		MockSnakeableNoRun: *NewMockSnakeableNoRun(),
		RunFunc:            func(customInterface) error { return nil },
	}

	ctx := context.Background()

	rootcmd := NewMockSnakeableNoRun()

	var err error
	ctx, err = NewRootCommand(ctx, rootcmd)
	if err != nil {
		t.Fatal(err)
	}

	ctx = RegisterBindingResolver(ctx, func(c context.Context) (context.Context, error) {
		return c, nil
	})

	ctx = RegisterBindingResolver(ctx, func(context.Context) (customInterface, error) {
		cms := customStruct{}
		return &cms, nil
	})

	ctx, err = NewCommand(ctx, "hello123", tt)
	assert.ErrorIs(t, err, tt.ExpectedNewCommandError())

	if err != nil {
		return
	}

	for _, b := range tt.Bindings() {
		ctx = Bind(ctx, reflect.ValueOf(b).Interface(), b)
	}

	os.Args = []string{"x", "hello123"}

	assm := Assemble(ctx)

	err = assm.ExecuteContext(ctx)

	assert.ErrorIs(t, err, tt.ExpectedRunCommandError())
}
func TestExpectInvalidContextResolver(t *testing.T) {

	testContextToBeInjected := context.Background()

	testContextToBeInjected = context.WithValue(testContextToBeInjected, "test", "test")

	tt := &SnakeableWithCustomInterfaceRunFuncContext{
		MockSnakeableNoRun: *NewMockSnakeableNoRun(),
		RunFunc: func(ctx context.Context, ci customInterface) error {

			if ctx.Value("test") == nil {
				return errors.New("context not passed")
			}

			return nil
		},
	}

	ctx := context.Background()

	rootcmd := NewMockSnakeableNoRun()

	var err error
	ctx, err = NewRootCommand(ctx, rootcmd)
	if err != nil {
		t.Fatal(err)
	}

	ctx = RegisterBindingResolver(ctx, func(cctx context.Context) (context.Context, error) {

		cctx = testContextToBeInjected

		cctx = Bind(cctx, (*customInterface)(nil), &customStruct{})

		return cctx, nil
	})

	ctx, err = NewCommand(ctx, "hello123", tt)
	assert.ErrorIs(t, err, tt.ExpectedNewCommandError())

	if err != nil {
		return
	}

	for _, b := range tt.Bindings() {
		ctx = Bind(ctx, reflect.ValueOf(b).Interface(), b)
	}

	os.Args = []string{"x", "hello123"}

	assm := Assemble(ctx)

	err = assm.ExecuteContext(context.TODO())

	assert.ErrorIs(t, err, ErrInvalidContextResolver)
}

func TestZeroLogProblemSimple(t *testing.T) {
	rootcmd := NewMockSnakeableNoRun()

	ctx := context.Background()

	ctx, err := NewRootCommand(ctx, rootcmd)
	if err != nil {
		t.Fatal(err)
	}

	zlWriterToBeInjected := zerolog.New(zerolog.NewTestWriter(nil))
	zlWriterToBeInjectedPtr := &zlWriterToBeInjected

	ctx = RegisterBindingResolver(ctx, func(cctx context.Context) (context.Context, error) {

		cctx = zlWriterToBeInjectedPtr.WithContext(cctx)

		zlg := zerolog.Ctx(cctx)

		assert.Equal(t, zlg, zlWriterToBeInjectedPtr)

		return cctx, nil
	})

	ctx = RegisterBindingResolver(ctx, func(cctx context.Context) (customInterface, error) {
		zlg := zerolog.Ctx(cctx)

		if zlg == nil {
			return nil, errors.New("zerolog not passed")
		}

		assert.Equal(t, zlWriterToBeInjectedPtr, zlg)

		cms := customStruct{}
		return &cms, nil
	})

	tt := &SnakeableWithCustomInterfaceRunFuncContext{
		MockSnakeableNoRun: *NewMockSnakeableNoRun(),
		RunFunc: func(ctx context.Context, ci customInterface) error {

			return nil
		},
	}

	ctx, err = NewCommand(ctx, "hello123", tt)
	if err != nil {
		t.Fatal(err)
	}
	os.Args = []string{"x", "hello123"}

	cmd := Assemble(ctx)

	err = cmd.ExecuteContext(context.TODO())
	assert.ErrorIs(t, err, tt.ExpectedRunCommandError())

}

func TestGetAlreadyBoundSimple(t *testing.T) {
	rootcmd := NewMockSnakeableNoRun()

	ctx := context.Background()

	ctx, err := NewRootCommand(ctx, rootcmd)
	if err != nil {
		t.Fatal(err)
	}

	type alreadyBound struct {
	}

	ctx = RegisterBindingResolver(ctx, func(cctx context.Context) (context.Context, error) {

		cctx = Bind(cctx, (*alreadyBound)(nil), &alreadyBound{})

		return cctx, nil
	})

	ctx = RegisterBindingResolver(ctx, func(cctx context.Context) (customInterface, error) {

		arb, ok := GetAlreadyBound[*alreadyBound](cctx)
		assert.True(t, ok)
		assert.NotNil(t, arb)

		cms := customStruct{}
		return &cms, nil
	})

	tt := &SnakeableWithCustomInterfaceRunFuncContext{
		MockSnakeableNoRun: *NewMockSnakeableNoRun(),
		RunFunc: func(ctx context.Context, ci customInterface) error {

			return nil
		},
	}

	ctx, err = NewCommand(ctx, "hello123", tt)
	if err != nil {
		t.Fatal(err)
	}
	os.Args = []string{"x", "hello123"}

	cmd := Assemble(ctx)

	err = cmd.ExecuteContext(context.TODO())
	assert.ErrorIs(t, err, tt.ExpectedRunCommandError())

}

func ptr[T any](t T) *T {
	return &t
}

func TestBindGBehavesAsExpected(t *testing.T) {

	// conver to table test
	tests := []struct {
		name string

		boundfn func(ctx context.Context) (any, bool)
		ctxfn   func(context.Context) context.Context
	}{
		{
			name: "string",
			boundfn: func(ctx context.Context) (any, bool) {
				return GetAlreadyBound[string](ctx)
			},
			ctxfn: func(ctx context.Context) context.Context {
				return BindG[string](ctx, "string")
			},
		},
		{
			name: "int",
			boundfn: func(ctx context.Context) (any, bool) {
				return GetAlreadyBound[int](ctx)
			},
			ctxfn: func(ctx context.Context) context.Context {
				return BindG[int](ctx, 1)
			},
		},
		{
			name: "uint64",
			boundfn: func(ctx context.Context) (any, bool) {
				return GetAlreadyBound[uint64](ctx)
			},
			ctxfn: func(ctx context.Context) context.Context {
				return BindG[uint64](ctx, 1)
			},
		},
		{
			name: "struct",
			boundfn: func(ctx context.Context) (any, bool) {
				return GetAlreadyBound[customStruct](ctx)
			},
			ctxfn: func(ctx context.Context) context.Context {
				return BindG[customStruct](ctx, customStruct{})
			},
		},
		{
			name: "struct ptr",
			boundfn: func(ctx context.Context) (any, bool) {
				return GetAlreadyBound[*customStruct](ctx)
			},
			ctxfn: func(ctx context.Context) context.Context {
				return BindG[*customStruct](ctx, &customStruct{})
			},
		},
		{
			name: "interface",
			boundfn: func(ctx context.Context) (any, bool) {
				return GetAlreadyBound[customInterface](ctx)
			},
			ctxfn: func(ctx context.Context) context.Context {
				return BindG[customInterface](ctx, &customStruct{})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gctx := context.Background()
			gctx = tt.ctxfn(gctx)
			bnd, ok := tt.boundfn(gctx)
			assert.True(t, ok)
			require.NotNil(t, bnd)

		})
	}

}
