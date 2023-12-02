package sample

import (
	"context"
	"io"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/walteh/snake"
	"github.com/walteh/snake/example/resolvers"
	"github.com/walteh/snake/scobra"
)

var _ scobra.SCobra = (*Handler)(nil)

type Handler struct {
	Value string `default:"default"`
	Cool  bool   `default:"false"`

	curr int

	args args
}

type args struct {
	Context context.Context
	Read    io.Reader
	Write   io.Writer
	Enum    resolvers.SampleEnum
	Br      io.ByteReader
	Bw      io.ByteWriter
	Bs      io.ByteScanner
}

func (me *Handler) Args() *args {
	return &me.args
}

func (*Handler) Name() string {
	return "sample"
}

func (*Handler) Description() string {
	return "sample description"
}

func (me *Handler) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use: me.Name(),
	}

	cmd.Args = cobra.ExactArgs(0)

	// cmd.Flags().StringVar(&me.Value, "value", "default", "value to print")
	cmd.Flags().BoolVar(&me.Cool, "cool", false, "cool value")

	return cmd
}

func (me *Handler) Run(
	ctx context.Context,
	read io.Reader,
	write io.Writer,
	en resolvers.SampleEnum,
	br io.ByteReader, bw io.ByteWriter, bs io.ByteScanner,
) (snake.Output, error) {
	me.args.Context = ctx
	me.args.Read = read
	me.args.Write = write
	me.args.Enum = en
	me.args.Br = br
	me.args.Bw = bw
	me.args.Bs = bs

	me.curr += 1
	return &snake.TableOutput{
		ColumnNames: []string{"name", "value"},
		RowValueData: [][]any{
			{"value", me.Value},
			{"cool", me.Cool},
			{"curr", me.curr},
		},
		RowValueColors: [][]*color.Color{
			{color.New(color.FgHiGreen), color.New(color.Bold)},
			{color.New(color.FgHiRed), color.New(color.Bold)},
			{color.New(color.FgHiBlue), color.New(color.FgBlack)},
		},
	}, nil

}
