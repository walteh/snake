package swails

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/walteh/snake/szerolog"
)

type ContextResolver struct {
}

func (me *ContextResolver) Run() (context.Context, error) {

	ctx := context.Background()

	ctx = szerolog.NewConsoleLoggerContext(ctx, zerolog.TraceLevel, os.Stdout)

	return ctx, nil
}
