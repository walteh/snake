package root

import (
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/walteh/snake"
)

var _ snake.Flagged = (*ContextResolver)(nil)

type DoubleResolver struct {
}

func (me *DoubleResolver) Flags(_ *pflag.FlagSet) {
}

func (me *DoubleResolver) Run(cmd *cobra.Command) (io.Reader, io.Writer, error) {

	return strings.NewReader("hello"), &strings.Builder{}, nil
}
