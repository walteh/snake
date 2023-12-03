package swails

import (
	"context"

	"github.com/walteh/snake"
)

var (
	_ snake.SnakeImplementation[SWails] = &WailsSnake{}
)

type WailsSnake struct {
	snake  snake.Snake
	binder *snake.Binder
	inputs map[string]snake.Input
}

type SWails interface {
	Name() string
	Description() string
}

func NewCommandResolver(s SWails) snake.TypedResolver[SWails] {
	return snake.MustGetTypedResolver(s)
}

func (me *WailsSnake) ManagedResolvers(_ context.Context) []snake.Resolver {
	return []snake.Resolver{}
}

func (me *WailsSnake) Decorate(ctx context.Context, self SWails, snk snake.Snake, inputs []snake.Input, mw []snake.Middleware) error {

	return nil
}

func (me *WailsSnake) OnSnakeInit(ctx context.Context, snk snake.Snake) error {

	me.snake = snk

	commands := me.snake.Resolvers()

	var inputs []snake.Input

	for _, cmd := range commands {
		snki, err := snake.InputsFor(cmd, me.snake.Enums()...)
		if err != nil {
			return err
		}

		inputs = append(inputs, snki...)
	}

	for _, input := range inputs {
		me.inputs[input.Name()] = input
	}

	return nil
}

func (me *WailsSnake) ResolveEnum(typ string, opts []string) (string, error) {
	return "", nil
}

func (me *WailsSnake) ProvideContextResolver() snake.Resolver {
	return snake.MustGetResolverFor[context.Context](&ContextResolver{})
}

func NewWailsSnake(ctx context.Context) *WailsSnake {

	me := &WailsSnake{
		binder: snake.NewBinder(),
		inputs: make(map[string]snake.Input),
	}

	return me
}

func ExecuteHandlingError(ctx context.Context, cmd *WailsSnake) {

}
