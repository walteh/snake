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
	A bool
	B bool
}

func (me *DoubleResolver) Flags(a *pflag.FlagSet) {
	a.BoolVarP(&me.A, "a", "a", false, "")
	a.BoolVarP(&me.B, "b", "b", false, "")
}

func (me *DoubleResolver) Run(cmd *cobra.Command) (io.Reader, io.Writer, error) {

	return strings.NewReader("hello"), &strings.Builder{}, nil
}
