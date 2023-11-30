package snake

import (
	"context"
	"encoding/json"
	"io"

	"github.com/fatih/color"
	"github.com/go-faster/errors"
)

type OutputHandler interface {
	HandleLongRunningOutput(ctx context.Context, out *LongRunningOutput) error
	HandleRawTextOutput(ctx context.Context, out *RawTextOutput) error
	HandleTableOutput(ctx context.Context, out *TableOutput) error
	HandleJSONOutput(ctx context.Context, out *JSONOutput) error
	HandleNilOutput(ctx context.Context, out *NilOutput) error
	HandleFileOutput(ctx context.Context, out *FileOutput) error
}

func (*LongRunningOutput) IsOutput() {}
func (*RawTextOutput) IsOutput()     {}
func (*TableOutput) IsOutput()       {}
func (*JSONOutput) IsOutput()        {}
func (*NilOutput) IsOutput()         {}
func (*FileOutput) IsOutput()        {}

type Output interface {
	IsOutput()
}

type FileOutput struct {
	Dir   string
	Mkdir bool
	Data  map[string]io.Reader
}

type LongRunningOutput struct {
	Start func(context.Context) error
}

type RawTextOutput struct {
	Data string
}

type TableOutput struct {
	ColumnNames    []string
	RowValueData   [][]any
	RowValueColors [][]*color.Color
}

type JSONOutput struct {
	Data json.RawMessage
}

type NilOutput struct{}

func HandleOutput(ctx context.Context, handler OutputHandler, out Output) error {
	if handler == nil {
		return errors.Errorf("trying to handle output with no handler provided - %T", out)
	}
	switch t := out.(type) {
	case *LongRunningOutput:
		return handler.HandleLongRunningOutput(ctx, t)
	case *RawTextOutput:
		return handler.HandleRawTextOutput(ctx, t)
	case *TableOutput:
		clength := len(t.ColumnNames)
		if len(t.RowValueData) != len(t.RowValueColors) {
			return errors.Errorf("table output data (%d) does not match colors (%d)", len(t.RowValueData), len(t.RowValueColors))
		}
		for _, row := range t.RowValueData {
			if len(row) != clength {
				return errors.Errorf("table output column names (%d) do not match data (%d)", clength, len(row))
			}
		}
		for _, row := range t.RowValueColors {
			if len(row) != clength {
				return errors.Errorf("table output column names (%d) do not match data (%d)", clength, len(row))
			}
		}
		return handler.HandleTableOutput(ctx, t)
	case *JSONOutput:
		return handler.HandleJSONOutput(ctx, t)
	case *NilOutput:
		return handler.HandleNilOutput(ctx, t)
	case *FileOutput:
		return handler.HandleFileOutput(ctx, t)
	default:
		return errors.Errorf("unknown output type %T", t)
	}
}
