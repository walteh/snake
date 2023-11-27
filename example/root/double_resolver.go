package root

import (
	"io"
	"strings"

	"github.com/spf13/cobra"
)

type DoubleResolver struct {
	A bool `usage:"A" default:"true"`
	B bool `usage:"B" default:"true"`
}

func (me *DoubleResolver) Run(cmd *cobra.Command) (io.Reader, io.Writer, error) {

	return strings.NewReader("hello"), &strings.Builder{}, nil
}
