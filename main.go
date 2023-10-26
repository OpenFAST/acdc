package main

import (
	"embed"

	"github.com/labstack/gommon/log"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	// "gonum.org/v1/gonum/blas/blas64"
	// "gonum.org/v1/gonum/lapack/lapack64"
	// blaslib "gonum.org/v1/netlib/blas/netlib"
	// lapacklib "gonum.org/v1/netlib/lapack/netlib"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	// blas64.Use(blaslib.Implementation{})
	// lapack64.Use(lapacklib.Implementation{})

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
