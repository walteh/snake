package root

import (
	"io"
	"strings"

	"github.com/spf13/cobra"
)

// var _ scobra.Flagged = (*ContextResolver)(nil)

type TripleResolver struct {
}

func (me *TripleResolver) Run(cmd *cobra.Command) (io.ByteReader, io.ByteWriter, io.ByteScanner, error) {
	return strings.NewReader("hello"), &strings.Builder{}, strings.NewReader("hello"), nil
}
