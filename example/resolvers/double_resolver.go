package resolvers

import (
	"io"
	"strings"
)

type DoubleResolver struct {
	A bool `usage:"A" default:"true"`
	B bool `usage:"B" default:"true"`
}

func (me *DoubleResolver) Run() (io.Reader, io.Writer, error) {

	return strings.NewReader("hello"), &strings.Builder{}, nil
}
