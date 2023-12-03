package main

import (
	"context"
	"embed"

	"github.com/walteh/snake"
	"github.com/walteh/snake/example/resolvers"
	"github.com/walteh/snake/example/root/basic"
	"github.com/walteh/snake/example/root/sample"
	"github.com/walteh/snake/swails"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	ctx := context.Background()

	swail := swails.NewWailsSnake(ctx)

	resolvers := resolvers.LoadResolvers()

	commands := []snake.Resolver{
		swails.NewCommandResolver(&basic.Handler{}),
		swails.NewCommandResolver(&sample.Handler{}),
	}

	// runtime.EventsOn(ctx, "wails:ready", func() {
	// 	// fmt.Println("Wails is ready!")
	// })

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "wails",
		Width:  1024,
		Height: 768,

		Debug: options.Debug{
			OpenInspectorOnStartup: true,
		},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},

		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			_, err := snake.NewSnakeWithOpts(ctx, swail, &snake.NewSnakeOpts{
				Resolvers: append(commands, resolvers...),
				OverrideEnumResolver: func(typ string, opts []string) (string, error) {
					return "y", nil
				},
			})

			if err != nil {
				panic(err)
			}
		},

		Bind: []interface{}{
			// app,
			swail,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
