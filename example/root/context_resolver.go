package root

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/walteh/snake/scobra"
)

var _ scobra.Flagged = (*ContextResolver)(nil)

type ContextResolver struct {
	Quiet   bool
	Debug   bool
	Version bool
}

func (me *ContextResolver) Flags(flgs *pflag.FlagSet) {
	flgs.BoolVarP(&me.Quiet, "quiet", "q", false, "Do not print any output")
	flgs.BoolVarP(&me.Debug, "debug", "d", false, "Print debug output")
	flgs.BoolVarP(&me.Version, "version", "v", false, "Print version and exit")
}

func (me *ContextResolver) Run(cmd *cobra.Command) (context.Context, error) {

	var level zerolog.Level
	if me.Debug {
		level = zerolog.TraceLevel
	} else if me.Quiet {
		level = zerolog.NoLevel
	} else {
		level = zerolog.InfoLevel
	}

	ctx := context.Background()

	ctx = zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger().Level(level).WithContext(ctx)

	return ctx, nil
}
