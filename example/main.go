package main

import (
	"context"
	"fmt"
	"os"

	"github.com/walteh/snake/example/root"
	"github.com/walteh/snake/scobra"
)

func main() {

	ctx := context.Background()

	_, cmd, _, err := root.NewCommand(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scobra.ExecuteHandlingError(ctx, cmd)
}
