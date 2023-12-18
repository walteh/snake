package basic

import (
	"fmt"

	"github.com/walteh/snake"
	"github.com/walteh/snake/example/resolvers"
)

func Runner() snake.Runner {
	return snake.GenRunCommand_In01_Out02(&Handler{})
}

type Handler struct {
	Value string `default:"default"`
}

func (me *Handler) Name() string {
	return "basic"
}

func (me *Handler) Description() string {
	return "basic description"
}

func (me *Handler) Run(dat resolvers.CustomInterface) (snake.Output, error) {
	return &snake.RawTextOutput{
		Data: fmt.Sprintf("hello %s, my value is %s", me.Value, dat),
	}, nil
}
