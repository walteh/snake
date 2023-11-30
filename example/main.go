package main

import (
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/walteh/snake/example/root"
	"github.com/walteh/snake/scobra"
)

func main() {

	ctx := context.Background()

	_, cmd, _, err := root.NewCommand(ctx)
	if err != nil {
		if !scobra.IsHandledByPrintingToConsole(err) {
			_, _ = fmt.Print(err)
		}
		os.Exit(1)
	}

	ctx = cmd.RootCommand.Context()

	str, err := scobra.DecorateTemplate(ctx, cmd.RootCommand, &scobra.DecorateOptions{
		Headings: color.New(color.FgCyan, color.Bold),
		ExecName: color.New(color.FgHiGreen, color.Bold),
		Commands: color.New(color.FgHiRed, color.Faint),
	})
	if err != nil {
		if !scobra.IsHandledByPrintingToConsole(err) {
			_, _ = fmt.Print(err)
		}
		os.Exit(1)
	}

	cmd.RootCommand.SetUsageTemplate(str)

	// cmd.SilenceErrors = true

	if err := cmd.RootCommand.ExecuteContext(ctx); err != nil {
		if !scobra.IsHandledByPrintingToConsole(err) {
			_, _ = fmt.Print(err)
		}
		os.Exit(1)
	}

	fmt.Println("done")

}
