package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/wailsapp/wails/v2/pkg/logger"
)

var LogFile = "log.txt"
var LogDir = filepath.Join(xdg.StateHome, "acdc")
var LogPath = filepath.Join(LogDir, LogFile)

func NewLogger() (logger.Logger, error) {
	if err := os.MkdirAll(LogDir, 0777); err != nil {
		return nil, fmt.Errorf("error creating log directory '%s': %w", LogDir, err)
	}
	return logger.NewFileLogger(LogPath), nil
}
