package sample

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/walteh/snake"
)

var _ snake.Cobrad = (*Handler)(nil)
var _ snake.Flagged = (*Handler)(nil)

type Handler struct {
	Value string
	Cool  bool
}

func (me *Handler) Flags(s *pflag.FlagSet) {
	s.StringVar(&me.Value, "value", "default", "value to print")
	s.BoolVar(&me.Cool, "cool", false, "cool value")
}

func (me *Handler) Cobra() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sample",
		Short: "run a server for retab code using the Language Server Protocol",
	}

	cmd.Args = cobra.ExactArgs(0)

	return cmd
}

func (me *Handler) ParseArguments(_ context.Context, _ *cobra.Command, _ []string) error {

	return nil

}

func (me *Handler) Run(ctx context.Context) error {
	return nil
	// return NewServe().Run(debug.WithInstance(ctx, "./de.bug", "serve"), nil)
}
