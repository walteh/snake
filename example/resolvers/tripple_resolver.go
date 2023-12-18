package resolvers

import (
	"io"
	"strings"

	"github.com/walteh/snake"
)

// var _ scobra.Flagged = (*ContextResolver)(nil)

func TripleRunner() snake.Runner {
	return snake.GenRunResolver_In00_Out04(&TripleResolver{})
}

type TripleResolver struct {
}

func (me *TripleResolver) Run() (io.ByteReader, io.ByteWriter, io.ByteScanner, error) {
	return strings.NewReader("hello"), &strings.Builder{}, strings.NewReader("hello"), nil
}
