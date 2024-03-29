package sample

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/walteh/snake"
	"github.com/walteh/snake/example/resolvers"
)

func (me *Handler) RegisterRunFunc() snake.RunFunc {
	return snake.GenRunCommand_In08_Out02(me)
}

type Handler struct {
	Value      string   `default:"default"`
	Cool       bool     `default:"false"`
	CoolArray  []string `default:"a,b,c"`
	CoolIntArr []int    `default:"1,2,3"`

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

func (me *Handler) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sample",
		Short: "sample command",
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
	br io.ByteReader, bw io.ByteWriter, bs io.ByteScanner, out snake.Stdout,
) (snake.Output, error) {
	me.args.Context = ctx
	me.args.Read = read
	me.args.Write = write
	me.args.Enum = en
	me.args.Br = br
	me.args.Bw = bw
	me.args.Bs = bs

	fmt.Fprintf(out, "value: %s\n", me.Value)

	fmt.Println("cool: ", me.Cool)

	me.curr += 1
	return &snake.TableOutput{
		ColumnNames: []string{"name", "value"},
		RowValueData: [][]any{
			{"value", me.Value},
			{"cool", me.Cool},
			{"curr", me.curr},
		},
		RowValueColors: [][]string{
			{"green", "red"}, {"blue", "black"}, {"yellow", "white"},
		},
		RawData: me,
	}, nil

}
