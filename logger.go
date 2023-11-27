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
	execPath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	LogDir = filepath.Dir(execPath)
	LogPath = filepath.Join(LogDir, "ACDC.log")
	// if err := os.MkdirAll(LogDir, 0777); err != nil {
	// 	return nil, fmt.Errorf("error creating log directory '%s': %w", LogDir, err)
	// }
	fmt.Println("writing log at: " + LogPath)
	return logger.NewFileLogger(LogPath), nil
}
