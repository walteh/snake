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

func (*Handler) Name() string {
	return "basic"
}

func (*Handler) Description() string {
	return "basic description"
}

func (*Handler) Image() string {
	return "https://tailwindui.com/img/logos/48x48/tuple.svg"
}

func (*Handler) Emoji() string {
	return "🚀"
}

func (*Handler) Command() *cobra.Command {
	return &cobra.Command{
		Use: "basic",
	}
}

func (me *Handler) Run(dat resolvers.DependantResolverString) (snake.Output, error) {
	return &snake.RawTextOutput{
		Data: fmt.Sprintf("hello %s, my value is %s", me.Value, dat),
	}, nil
}
