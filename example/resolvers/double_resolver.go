package resolvers

import (
	"io"
	"strings"

	"github.com/walteh/snake"
)

func DoubleRunner() snake.Runner {
	return snake.GenRunResolver_In00_Out04(&TripleResolver{})
}

type DoubleResolver struct {
	A bool `usage:"A" default:"true"`
	B bool `usage:"B" default:"true"`
}

func (me *DoubleResolver) Run() (io.Reader, io.Writer, error) {

	return strings.NewReader("hello"), &strings.Builder{}, nil
}
