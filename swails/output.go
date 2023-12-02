package swails

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/walteh/snake"
)

// YO could we have more inputs for the ouputs? maybe add some specific flags to them if they have a run method?
// like being able to pass a --json flag to the table output and it will convert it to json? or a csv flag, or a file flag.

var _ snake.OutputHandler = (*OutputHandler)(nil)

type OutputHandler struct {
	output *WailsHTMLResponse
	stdio  io.Writer
}

func NewOutputHandler(cmd io.Writer) *OutputHandler {
	return &OutputHandler{
		stdio: cmd,
	}
}

func (me *OutputHandler) HandleJSONOutput(ctx context.Context, out *snake.JSONOutput) error {

	// Convert the output data to JSON format
	jsonData, err := json.MarshalIndent(out.Data, "", "\t")
	if err != nil {
		return err // Handle or return the error appropriately
	}

	// Print the formatted JSON to the command's output
	_, _ = fmt.Fprintln(me.stdio, string(jsonData))

	tmpl := `
	<div>
					<code class="language-json">
						%s
					</code>
	</div>
	`

	me.output = &WailsHTMLResponse{
		HTML: fmt.Sprintf(tmpl, string(jsonData)),
	}

	return nil
}

// HandleLongRunningOutput implements sbind.OutputHandler.
func (*OutputHandler) HandleLongRunningOutput(ctx context.Context, out *snake.LongRunningOutput) error {
	return out.Start(ctx)
}

// HandleRawTextOutput implements sbind.OutputHandler.
func (me *OutputHandler) HandleRawTextOutput(ctx context.Context, out *snake.RawTextOutput) error {

	tmpl := `
	<div>
					<code class="language-txt">
						%s
					</code>
	</div>
	`

	me.output = &WailsHTMLResponse{
		HTML: fmt.Sprintf(tmpl, out.Data),
	}

	return nil
}

// HandleTableOutput implements sbind.OutputHandler.
func (me *OutputHandler) HandleTableOutput(ctx context.Context, out *snake.TableOutput) error {

	tmpl := `
	<div>
		<table class="table table-striped table-bordered table-hover">
			<thead>
				<tr>
					%s
				</tr>
			</thead>
			<tbody>
				%s
			</tbody>
		</table>
	</div>
	`

	header := ""
	for _, col := range out.ColumnNames {
		header += fmt.Sprintf("<th>%s</th>", col)
	}

	rows := ""
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
		rows += fmt.Sprintf("<tr><td>%s</td></tr>", strings.Join(strdat, "</td><td>"))
	}

	me.output = &WailsHTMLResponse{
		HTML: fmt.Sprintf(tmpl, header, rows),
	}

	return nil
}

// HandleNilOutput implements sbind.OutputHandler.
func (me *OutputHandler) HandleNilOutput(ctx context.Context, out *snake.NilOutput) error {
	me.output = &WailsHTMLResponse{
		HTML: "nil output",
	}
	return nil
}

// HandleFileOutput implements sbind.OutputHandler.
func (me *OutputHandler) HandleFileOutput(ctx context.Context, out *snake.FileOutput) error {
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
		HTML: string(res),
	}

	return nil
}
