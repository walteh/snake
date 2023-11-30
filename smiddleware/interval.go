package smiddleware

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/walteh/snake"
)

type IntervalMiddleware struct {
	Interval     time.Duration `usage:"interval between runs"`
	NumberOfRuns int           `usage:"number of runs (-1 for infinite)"`
}

func (me *IntervalMiddleware) Wrap(base snake.MiddlewareFunc) snake.MiddlewareFunc {
	return func(ctx context.Context) error {

		start := time.Now()

		ctx = zerolog.Ctx(ctx).With().
			Dur("interval", me.Interval).
			Str("middleware", "IntervalMiddleware").
			Logger().
			WithContext(ctx)

		zerolog.Ctx(ctx).Trace().Msg("starting")

		times := 0

		defer func() {
			zerolog.Ctx(ctx).Trace().
				Int("runs", times).
				Dur("run_time", time.Since(start)).
				Msgf("done")
		}()

		runner := func(ctx context.Context) error {
			times++

			zerolog.Ctx(ctx).Trace().
				Int("run_number", times).
				Dur("run_time", time.Since(start)).
				Msgf("running")

			if err := base(ctx); err != nil {
				return err
			}
			return nil
		}

		if me.NumberOfRuns != 0 {
			err := runner(ctx)
			if err != nil {
				return err
			}
		}

		for me.NumberOfRuns == -1 || times < me.NumberOfRuns {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(me.Interval):
				err := runner(ctx)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}
}

func NewIntervalMiddleware() *IntervalMiddleware {
	return &IntervalMiddleware{
		Interval:     time.Second * 1,
		NumberOfRuns: -1,
	}
}

func NewIntervalMiddlewareWithDefault(def time.Duration) *IntervalMiddleware {
	return &IntervalMiddleware{
		Interval:     def,
		NumberOfRuns: -1,
	}
}
