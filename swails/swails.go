package swails

import (
	"context"

	"github.com/walteh/snake"
)

var (
	_ snake.SnakeImplementationTyped[SWails] = &WailsSnake{}
)

type WailsSnake struct {
	snake            snake.Snake
	binder           *snake.Binder
	inputs           map[string]map[string]snake.Input
	emitter          WailsEmitter
	lifecycleContext context.Context
}

type SWails interface {
	Name() string
	Description() string
	Image() string
	Emoji() string
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

	for _, cmd := range commands {
		snki, err := snake.InputsFor(cmd, me.snake.Enums()...)
		if err != nil {
			return err
		}
		for _, input := range snki {

			if input.Shared() {
				if _, ok := me.inputs["shared"]; !ok {
					me.inputs["shared"] = make(map[string]snake.Input)
				}
				me.inputs["shared"][input.Name()] = input
			} else {
				if typ, ok := cmd.(snake.TypedResolver[SWails]); ok {
					if _, ok := me.inputs[typ.TypedRef().Name()]; !ok {
						me.inputs[typ.TypedRef().Name()] = make(map[string]snake.Input)
					}
					me.inputs[typ.TypedRef().Name()][input.Name()] = input
				}

			}
		}

	}

	return nil
}

func inputName(cmd SWails, input snake.Input) string {
	if input.Shared() {
		return input.Name()
	}
	return cmd.Name() + input.Name()
}

func (me *WailsSnake) ResolveEnum(typ string, opts []string) (string, error) {
	return "", nil
}

func (me *WailsSnake) ProvideContextResolver() snake.Resolver {
	return snake.MustGetResolverFor[context.Context](&ContextResolver{})
}

func NewWailsSnake(ctx context.Context, emitter WailsEmitter) *WailsSnake {

	me := &WailsSnake{
		binder:  snake.NewBinder(),
		inputs:  make(map[string]map[string]snake.Input),
		emitter: emitter,
	}

	return me
}

func (me *WailsSnake) SetLifecycleContext(ctx context.Context) {
	me.lifecycleContext = ctx
}

func ExecuteHandlingError(ctx context.Context, cmd *WailsSnake) {

}
