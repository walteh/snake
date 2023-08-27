package snake

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/k0kubun/pp/v3"
	"github.com/rs/zerolog"
)

func Ctx(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}

func NewJsonLogger() *zerolog.Logger {
	logger := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
	return &logger
}

func NewJsonStdErrLogger() *zerolog.Logger {
	logger := zerolog.New(os.Stderr).With().Timestamp().Caller().Logger()
	return &logger
}

func NewVerboseLoggerContext(ctx context.Context) context.Context {
	verboseLogger := NewVerboseConsoleLogger()
	return verboseLogger.WithContext(ctx)
}

func NewJsonLoggerContext(ctx context.Context) context.Context {
	jsonLogger := NewJsonLogger()
	return jsonLogger.WithContext(ctx)
}

func NewJsonStdErrLoggerContext(ctx context.Context) context.Context {
	jsonLogger := NewJsonStdErrLogger()
	return jsonLogger.WithContext(ctx)
}

func NewVerboseConsoleLogger() *zerolog.Logger {

	consoleOutput := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.StampMicro, NoColor: false}

	pretty := pp.New()

	pretty.SetColorScheme(pp.ColorScheme{})

	prettyerr := pp.New()
	prettyerr.SetExportedOnly(false)

	consoleOutput.FormatFieldValue = func(i interface{}) string {

		switch i := i.(type) {
		case error:
			return prettyerr.Sprint(i)
		// case string:
		// 	return i
		// case int:
		// 	return strconv.Itoa(i)
		// case int64:
		// 	return strconv.FormatInt(i, 10)
		// case float64:
		// 	return strconv.FormatFloat(i, 'f', -1, 64)
		case []byte:
			var g interface{}
			err := json.Unmarshal(i, &g)
			if err != nil {
				return pretty.Sprint(string(i))
			} else {
				return pretty.Sprint(g)
			}
		}

		return pretty.Sprint(i)
	}

	consoleOutput.FormatTimestamp = func(i interface{}) string {
		return time.Now().Format("[0000-00-00 | 15:04:05.000000]")
	}

	callerTrier := 0

	consoleOutput.FormatCaller = func(i interface{}) string {
		a := i.(string)
		tot := strings.Split(a, "/")
		if len(tot) == 3 {
			num := strings.Split(tot[2], ":")
			lll := 5 + len(num[0]) + len(num[1]) + len(tot[0]) + len(tot[1])
			if lll >= callerTrier {
				callerTrier = lll + 2
			}
			padding := strings.Repeat(" ", callerTrier-lll)
			return fmt.Sprintf("[%s:%s] %s:%s%s", tot[0], tot[1], color.BlueString(num[0]), color.New(color.FgBlue, color.Bold).Sprint(num[1]), padding)
		}

		// return the caller in blue in the console
		// make it black
		return fmt.Sprintf("\x1b[0m\x1b[34;1m%s\x1b[0m", i)
	}

	consoleOutput.PartsOrder = []string{"level", "time", "caller", "message"}

	consoleOutput.FieldsExclude = []string{"handler", "tags"}

	l := zerolog.New(consoleOutput).With().Caller().Timestamp().Logger()

	return &l

}
