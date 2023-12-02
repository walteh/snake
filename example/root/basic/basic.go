package basic

import (
	"fmt"

	"github.com/walteh/snake"
	"github.com/walteh/snake/example/resolvers"
)

type Handler struct {
	Value string `default:"default"`
}

func (*Handler) Name() string {
	return "basic"
}

func (*Handler) Description() string {
	return "basic description"
}

func (me *Handler) Run(dat resolvers.DependantResolverString) (snake.Output, error) {
	return &snake.RawTextOutput{
		Data: fmt.Sprintf("hello %s, my value is %s", me.Value, dat),
	}, nil
}
