package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const ProjectFile = "project.json"

// App struct
type App struct {
	ctx     context.Context
	Project *Project
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) OpenProjectDialog() (*Project, error) {

	// Open dialog so user can select the file
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
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

	// Open project
	return a.OpenProject(path)
}

// OpenProject opens the project at the given path
func (a *App) OpenProject(path string) (*Project, error) {

	// Load project
	var err error
	a.Project, err = LoadProject(path)
	if err != nil {
		return nil, err
	}

	// Set project path
	a.Project.Info.Path = path

	// set window title
	runtime.WindowSetTitle(a.ctx, "ACDC - "+path)

	return a.Project, nil
}

// SaveProject saves changes to the current project
func (a *App) SaveProject(path string) (*Project, error) {

	// Save project
	if _, err := a.Project.Save(path); err != nil {
		return nil, err
	}

	// Initialize return project
	p := &Project{Info: a.Project.Info}

	return p, nil
}

// SaveProjectDialog opens a dialog to select where to save the project
func (a *App) SaveProjectDialog() (*Project, error) {

	// Open dialog so user can select the file
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Save Project As",
		DefaultFilename: "project.json",
		Filters: []runtime.FileFilter{
			{DisplayName: "Projects (*.json)", Pattern: "*.json"},
		},
		CanCreateDirectories: true,
	})
	if err != nil {
		return nil, fmt.Errorf("error selecting project directory: %w", err)
	}

	// Create path
	if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		return nil, fmt.Errorf("error creating project directory '%s': %w", path, err)
	}

	// If project not loaded, create new project
	if a.Project == nil {
		a.Project = NewProject()
	}

	// Save project
	return a.SaveProject(path)
}

// SelectExec opens a dialog for the user to select an OpenFAST executable.
func (a *App) SelectExec() (*Project, error) {

	// Open dialog so user can select the executable path
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:                "Select OpenFAST Executable",
		CanCreateDirectories: false,
	})
	if err != nil {
		return nil, err
	}

	// If no path selected, return current exec
	if path == "" {
		return &Project{Info: a.Project.Info}, nil
	}

	// Get output from running the command
	output, _ := exec.Command(path).CombinedOutput()

	// If OpenFAST isn't in the output, return error
	if !bytes.Contains(output, []byte("OpenFAST")) {
		return nil, fmt.Errorf("'%s' is not an OpenFAST executable", path)
	}

	// Trim output
	if index := bytes.Index(output, []byte("Execution Info:")); index > -1 {
		output = output[:index-1]
	}

	// Update executable info in project
	a.Project.Exec = Exec{
		Path:    path,
		Version: string(bytes.TrimSpace(output)),
		Valid:   true,
	}

	// Save project
	if _, err := a.SaveProject(a.Project.Info.Path); err != nil {
		runtime.LogErrorf(a.ctx, "SelectExec: error saving project: %s", err)
	}

	// Initialize return project
	p := &Project{Info: a.Project.Info, Exec: a.Project.Exec}

	// Save project and return path
	return p, nil
}

func (a *App) ImportModelDialog() (*Project, error) {

	// Open dialog so user can select the file
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Open Model",
		Filters: []runtime.FileFilter{
			{DisplayName: "OpenFAST Model (*.fst)", Pattern: "*.fst"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error selecting OpenFAST file: %w", err)
	}

	// Parse model
	a.Project.Model, err = ParseModel(path)
	if err != nil {
		return nil, fmt.Errorf("error importing OpenFAST model '%s': %w", path, err)
	}

	// Save project
	if _, err := a.SaveProject(a.Project.Info.Path); err != nil {
		runtime.LogErrorf(a.ctx, "SelectExec: error saving project: %s", err)
	}

	// Initialize return Project
	p := &Project{Info: a.Project.Info, Model: a.Project.Model}

	// Parse and return model
	return p, nil
}

// UpdateModel saves changes to the current project
func (a *App) UpdateModel(model *Model) (*Project, error) {

	// Update model in project
	a.Project.Model = model

	// Save project
	if _, err := a.Project.Save(a.Project.Info.Path); err != nil {
		return nil, err
	}

	// Initialize return project
	p := &Project{Info: a.Project.Info}

	return p, nil
}
