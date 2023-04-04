package main

import (
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	// Create an instance of the app structure
	app := NewApp()

	// Create an instance of the project structure
	project := NewProject()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "ACDC",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			project.startup(ctx)
		},
		Bind: []interface{}{
			app,
			project,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
