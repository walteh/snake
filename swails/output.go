package swails

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"

	"github.com/walteh/snake"
)

// YO could we have more inputs for the ouputs? maybe add some specific flags to them if they have a run method?
// like being able to pass a --json flag to the table output and it will convert it to json? or a csv flag, or a file flag.

var _ snake.OutputHandler = (*OutputHandler)(nil)

func (me *OutputHandler) Stderr() io.Writer {
	return me.stdio
}

func (me *OutputHandler) Stdout() io.Writer {
	return me.stdio
}

func (me *OutputHandler) Stdin() io.Reader {
	return os.Stdin
}

type OutputHandler struct {
	output *WailsHTMLResponse
	stdio  io.Writer
}

func NewOutputHandler(cmd io.Writer) *OutputHandler {
	return &OutputHandler{
		stdio: cmd,
	}
}

func (me *OutputHandler) HandleJSONOutput(ctx context.Context, _ snake.Chan, out *snake.JSONOutput) error {

	// Convert the output data to JSON format
	jsonData, err := json.MarshalIndent(out.Data, "", "\t")
	if err != nil {
		return err // Handle or return the error appropriately
	}

	dat := make(map[string]any)

	err = json.Unmarshal(jsonData, &dat)
	if err != nil {
		return err
	}

	// Print the formatted JSON to the command's output
	_, _ = fmt.Fprintln(me.stdio, string(jsonData))

	me.output = &WailsHTMLResponse{
		Default: "json",
		JSON:    dat,
	}

	return nil
}

// HandleLongRunningOutput implements sbind.OutputHandler.
func (*OutputHandler) HandleLongRunningOutput(ctx context.Context, _ snake.Chan, out *snake.LongRunningOutput) error {
	return out.Start(ctx)
}

// HandleRawTextOutput implements sbind.OutputHandler.
func (me *OutputHandler) HandleRawTextOutput(ctx context.Context, _ snake.Chan, out *snake.RawTextOutput) error {

	me.output = &WailsHTMLResponse{
		Text:    out.Data,
		Default: "text",
	}

	return nil
}

// HandleTableOutput implements sbind.OutputHandler.
func (me *OutputHandler) HandleTableOutput(ctx context.Context, _ snake.Chan, out *snake.TableOutput) error {

	jsond := make([]map[string]any, len(out.RowValueData))

	data := [][]string{out.ColumnNames}
	stypes := [][]string{out.ColumnNames}

	for i, row := range out.RowValueData {
		strdat := make([]string, len(row))
		strsty := make([]string, len(row))
		for j, v := range row {

			strsty[j] = fmt.Sprintf("%s", out.RowValueColors[i][j])

			if reflect.TypeOf(v).Kind() == reflect.Ptr {
				v = reflect.ValueOf(v).Elem().Interface()
			}
			if v == nil {
				strdat[j] = "NULL"
				continue
			}
			strdat[j] = fmt.Sprintf("%v", v)
		}

		jsond[i] = make(map[string]interface{})
		for j, v := range row {
			jsond[i][out.ColumnNames[j]] = v
		}

		data = append(data, strdat)
		stypes = append(stypes, strsty)
	}

	me.output = &WailsHTMLResponse{
		Default:     "table",
		Table:       data,
		TableStyles: stypes,
		JSON:        jsond,
	}

	return nil
}

// HandleNilOutput implements sbind.OutputHandler.
func (me *OutputHandler) HandleNilOutput(ctx context.Context, _ snake.Chan, out *snake.NilOutput) error {
	me.output = &WailsHTMLResponse{
		Text: "no output",
	}
	return nil
}

// HandleFileOutput implements sbind.OutputHandler.
func (me *OutputHandler) HandleFileOutput(ctx context.Context, _ snake.Chan, out *snake.FileOutput) error {
	dir := out.Dir

	if dir == "" {
		dir = "."
	}

	dir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	r, writ := io.Pipe()

	if out.Mkdir {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	_, _ = fmt.Fprintln(writ, "")

	_, _ = fmt.Fprintf(writ, "writing %d files to %s\n", len(out.Data), dir)

	for name, content := range out.Data {
		dat, err := io.ReadAll(content)
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintf(writ, "writing %d bytes to %s...", len(dat), name)
		err = os.WriteFile(filepath.Join(dir, name), dat, 0644)
		if err != nil {
			_, _ = fmt.Fprintln(writ, "...failed")
			return err
		}
		_, _ = fmt.Fprintln(writ, "...done")
	}

	_, _ = fmt.Fprintln(writ, "done writing files")

	_, _ = fmt.Fprintln(writ, "")

	_ = writ.Close()

	res, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	me.output = &WailsHTMLResponse{
		Text: string(res),
	}

	return nil
}
