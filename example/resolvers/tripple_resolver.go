package resolvers

import (
	"io"
	"strings"
)

// var _ scobra.Flagged = (*ContextResolver)(nil)

type TripleResolver struct {
}

func (me *TripleResolver) Run() (io.ByteReader, io.ByteWriter, io.ByteScanner, error) {
	return strings.NewReader("hello"), &strings.Builder{}, strings.NewReader("hello"), nil
}
