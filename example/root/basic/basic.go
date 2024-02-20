package basic

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/walteh/snake"
	"github.com/walteh/snake/example/resolvers"
)

type Handler struct {
	Value string `default:"default"`
}

func (me *Handler) CobraCommand() *cobra.Command {
	return &cobra.Command{Use: "basic"}
}
func (me *Handler) RegisterRunFunc() snake.RunFunc {
	return snake.GenRunCommand_In01_Out02(me)
}

func (me *Handler) Run(dat resolvers.CustomInterface) (snake.Output, error) {
	return &snake.RawTextOutput{
		Data: fmt.Sprintf("hello %s, my value is %s", me.Value, dat),
	}, nil
}
