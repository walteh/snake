package sample

import (
	"context"
	"errors"
	"io"

	"github.com/spf13/cobra"
	"github.com/walteh/snake/scobra"
)

var _ scobra.SCobra = (*Handler)(nil)

type Handler struct {
	Value string
	Cool  bool
}

func (me *Handler) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sample",
		Short: "run a server for retab code using the Language Server Protocol",
	}

	cmd.Args = cobra.ExactArgs(0)

	cmd.Flags().StringVar(&me.Value, "value", "default", "value to print")
	cmd.Flags().BoolVar(&me.Cool, "cool", false, "cool value")

	return cmd
}

func (me *Handler) ParseArguments(_ context.Context, _ *cobra.Command, _ []string) error {

	return nil

}

func (me *Handler) Run(
	ctx context.Context,
	cmd *cobra.Command,
	arr []string,
	read io.Reader,
	write io.Writer,
	br io.ByteReader, bw io.ByteWriter, bs io.ByteScanner,
) error {
	arrs := []any{ctx, cmd, arr, read, write, br, bw, bs}
	for _, a := range arrs {
		if a == nil {
			return errors.New("something is nil")
		}
	}
	return nil
	// return NewServe().Run(debug.WithInstance(ctx, "./de.bug", "serve"), nil)
}
