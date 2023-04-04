package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const ProjectFile = "project.json"

type Info struct {
	Name string    `json:"Name"`
	Date time.Time `json:"Date"`
	Path string    `json:"Path"`
}

type Project struct {
	Info

	ExecPath string

	ctx context.Context
}

func NewProject() *Project {
	return &Project{}
}

func (p *Project) startup(ctx context.Context) {
	p.ctx = ctx
}

func (p *Project) Create() (*Info, error) {

	// Open dialog so user can select the file
	path, err := runtime.SaveFileDialog(p.ctx, runtime.SaveDialogOptions{
		Title:           "Create Project",
		DefaultFilename: "project.json",
		Filters: []runtime.FileFilter{
			{DisplayName: "Projects (*.json)", Pattern: "*.json"},
		},
		CanCreateDirectories: true,
	})
	if err != nil {
		return nil, fmt.Errorf("error selecting project directory: %w", err)
	}

	// If path not selected, return
	if path == "" {
		return nil, fmt.Errorf("no directory selected")
	}

	// Create path
	err = os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return nil, fmt.Errorf("error creating project directory '%s': %w", path, err)
	}

	// Set project path
	p.Path = path

	// Save project
	if err := p.Save(); err != nil {
		return nil, err
	}

	// Return info
	return &p.Info, nil
}

func (p *Project) Open() (*Info, error) {

	// Open dialog so user can select the file
	path, err := runtime.OpenFileDialog(p.ctx, runtime.OpenDialogOptions{
		Title: "Open Project",
		Filters: []runtime.FileFilter{
			{DisplayName: "Projects (*.json)", Pattern: "*.json"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error selecting project file: %w", err)
	}

	// If path not selected, return
	if path == "" {
		return nil, fmt.Errorf("no file selected")
	}

	// Read project file
	bs, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading project: %w", err)
	}

	// Parse project file
	if err := json.Unmarshal(bs, p); err != nil {
		return nil, fmt.Errorf("error parsing project: %w", err)
	}

	// Set project directory
	p.Path = path

	return &p.Info, nil
}

func (p *Project) Save() error {

	// Update time
	p.Date = time.Now()

	// Convert project to json
	bs, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		return fmt.Errorf("error encoding project: %w", err)
	}

	// Write project file
	err = os.WriteFile(p.Path, bs, 0777)
	if err != nil {
		return fmt.Errorf("error reading project: %w", err)
	}

	return nil
}

// SelectExec opens a dialog for the user to select an OpenFAST executable.
func (p *Project) SelectExec() (string, error) {

	// Open dialog so user can select the file
	path, err := runtime.OpenFileDialog(p.ctx, runtime.OpenDialogOptions{
		Title:                "Select OpenFAST Executable",
		CanCreateDirectories: false,
	})
	if err != nil {
		return path, err
	}

	// TODO: verify that executable works

	// Save path to project
	p.ExecPath = path

	// Save project and return path
	return path, p.Save()
}
