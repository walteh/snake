package sbind

import (
	"context"
	"encoding/json"

	"github.com/fatih/color"
)

type Output interface {
	Distribute(context.Context) error
}

type LongRunningOutput struct {
	Start func(context.Context) error
}

type RawTextOutput struct {
	Data string
}

type TableOutput struct {
	Data   [][]string
	Colors [][]color.Color
}

type JSONOutput struct {
	Data json.RawMessage
}
