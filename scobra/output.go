package scobra

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"reflect"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/walteh/snake/sbind"
)

// YO could we have more inputs for the ouputs? maybe add some specific flags to them if they have a run method?
// like being able to pass a --json flag to the table output and it will convert it to json? or a csv flag, or a file flag.

var _ sbind.OutputHandler = (*OutputHandler)(nil)

type OutputHandler struct {
	cmd *cobra.Command
}

func NewOutputHandler(cmd *cobra.Command) *OutputHandler {
	return &OutputHandler{
		cmd: cmd,
	}
}

// HandleJSONOutput implements sbind.OutputHandler.
func (*OutputHandler) HandleJSONOutput(ctx context.Context, out *sbind.JSONOutput) error {
	panic("unimplemented")
}

// HandleLongRunningOutput implements sbind.OutputHandler.
func (*OutputHandler) HandleLongRunningOutput(ctx context.Context, out *sbind.LongRunningOutput) error {
	return out.Start(ctx)
}

// HandleRawTextOutput implements sbind.OutputHandler.
func (me *OutputHandler) HandleRawTextOutput(ctx context.Context, out *sbind.RawTextOutput) error {
	me.cmd.Println("")

	me.cmd.Println(out.Data)

	me.cmd.Println("")
	return nil
}

// HandleTableOutput implements sbind.OutputHandler.
func (me *OutputHandler) HandleTableOutput(ctx context.Context, out *sbind.TableOutput) error {
	table := tablewriter.NewWriter(me.cmd.OutOrStdout())

	table.SetHeader(out.ColumnNames)

	for i, row := range out.RowValueData {

		strdat := make([]string, len(row))
		for j, v := range row {
			if reflect.TypeOf(v).Kind() == reflect.Ptr {
				v = reflect.ValueOf(v).Elem().Interface()
			}
			if v == nil {
				strdat[j] = out.RowValueColors[i][j].Sprint("NULL")
				continue
			}
			strdat[j] = out.RowValueColors[i][j].Sprintf("%v", v)
		}

		table.Append(strdat)
	}

	table.Render()

	return nil
}

// HandleNilOutput implements sbind.OutputHandler.
func (me *OutputHandler) HandleNilOutput(ctx context.Context, out *sbind.NilOutput) error {
	me.cmd.Println("nil output")
	return nil
}

// HandleFileOutput implements sbind.OutputHandler.
func (me *OutputHandler) HandleFileOutput(ctx context.Context, out *sbind.FileOutput) error {
	dir := out.Dir

	if dir == "" {
		dir = "."
	}

	dir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	if out.Mkdir {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	me.cmd.Println("")

	me.cmd.Printf("writing %d files to %s\n", len(out.Data), dir)

	for name, content := range out.Data {
		dat, err := io.ReadAll(content)
		if err != nil {
			return err
		}
		me.cmd.Printf("writing %d bytes to %s...", len(dat), name)
		err = os.WriteFile(filepath.Join(dir, name), dat, 0644)
		if err != nil {
			me.cmd.Println("...failed")
			return err
		}
		me.cmd.Println("...done")
	}

	me.cmd.Println("done writing files")

	me.cmd.Println("")

	return nil
}
