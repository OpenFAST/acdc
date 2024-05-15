package main

import (
	"embed"
	"flag"
	"fmt"

	"github.com/carlmjohnson/versioninfo"
	"github.com/labstack/gommon/log"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

var version = "0.4.0-" + versioninfo.Short()

func main() {

	showVersion := flag.Bool("version", false, "display version information")
	flag.Parse()
	if *showVersion {
		fmt.Println(version)
		return
	}

	// Create an instance of the app structure
	app := NewApp()

	// Create logger
	appLogger, err := NewLogger()
	if err != nil {
		log.Fatal(err)
	}

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "ACDC",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour:   &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:          app.startup,
		Bind:               []interface{}{app},
		Logger:             appLogger,
		LogLevel:           logger.TRACE,
		LogLevelProduction: logger.DEBUG,
	})
	if err != nil {
		log.Fatal(err)
	}
}
