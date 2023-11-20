package root

import (
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/walteh/snake"
)

var _ snake.Flagged = (*ContextResolver)(nil)

type TripleResolver struct {
}

func (me *TripleResolver) Flags(_ *pflag.FlagSet) {
}

func (me *TripleResolver) Run(cmd *cobra.Command) (io.ByteReader, io.ByteWriter, io.ByteScanner, error) {
	return strings.NewReader("hello"), &strings.Builder{}, strings.NewReader("hello"), nil
}
