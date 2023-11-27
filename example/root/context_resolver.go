package root

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

type ContextResolver struct {
	Quiet   bool `usage:"Do not print any output" default:"false"`
	Debug   bool `usage:"Print debug output" default:"false"`
	Version bool `usage:"Print version and exit" default:"false"`
	Cool    string
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
