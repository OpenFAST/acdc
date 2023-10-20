package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/sync/errgroup"
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

//------------------------------------------------------------------------------
// Project
//------------------------------------------------------------------------------

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

//------------------------------------------------------------------------------
// Turbine
//------------------------------------------------------------------------------

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

	// Parse model files
	files, err := ParseFiles(path)
	if err != nil {
		return nil, fmt.Errorf("error importing OpenFAST model '%s': %w", path, err)
	}

	// Get imported paths
	paths := []string{}
	for p := range files.PathMap {
		p, _ := filepath.Rel(filepath.Dir(path), p)
		paths = append(paths, p)
	}
	sort.Strings(paths)

	// Initialize models structure
	a.Project.Model = &Model{
		HasAero: ((len(files.AeroDyn)+len(files.AeroDyn14)) > 0 &&
			len(files.InflowWind) > 0),
		Files:         files,
		ImportedPaths: paths,
		Notes:         []string{},
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

//------------------------------------------------------------------------------
// Analysis
//------------------------------------------------------------------------------

// UpdateAnalysis runs Analysis.Calculate and saves results
func (a *App) UpdateAnalysis(analysis *Analysis) (*Project, error) {

	// Calculate analysis values
	if err := analysis.Calculate(); err != nil {
		return nil, err
	}

	// Update analysis in project
	a.Project.Analysis = analysis

	// Save project
	if _, err := a.Project.Save(a.Project.Info.Path); err != nil {
		return nil, err
	}

	// Initialize return project
	p := &Project{Info: a.Project.Info, Analysis: a.Project.Analysis}

	return p, nil
}

func (a *App) AddAnalysisCase() (*Project, error) {

	// Add new case to analysis
	a.Project.Analysis.Cases = append(a.Project.Analysis.Cases, NewCase())

	// Update analysis
	return a.UpdateAnalysis(a.Project.Analysis)
}

func (a *App) RemoveAnalysisCase(ID int) (*Project, error) {

	// Filter out case that matches ID
	tmp := []Case{}
	for _, c := range a.Project.Analysis.Cases {
		if c.ID != ID {
			tmp = append(tmp, c)
		}
	}
	a.Project.Analysis.Cases = tmp

	// Update analysis
	return a.UpdateAnalysis(a.Project.Analysis)
}

//------------------------------------------------------------------------------
// Evaluate
//------------------------------------------------------------------------------

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
	if index := bytes.LastIndex(output, []byte("****")); index > -1 {
		output = output[index+4:]
	}

	// Update executable info in project
	a.Project.Exec = &Exec{
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

func (a *App) EvaluateLinearization(c *Case, numCPUs int) ([]EvalStatus, error) {

	// Call existing cancel func
	EvalCancel(fmt.Errorf("new evaluation started"))

	// Create path to case directory
	caseDir := filepath.Join(strings.TrimSuffix(a.Project.Info.Path, filepath.Ext(a.Project.Info.Path)), fmt.Sprintf("case%02d", c.ID))
	if err := os.MkdirAll(caseDir, 0777); err != nil {
		return nil, fmt.Errorf("error creating directory '%s': %w", caseDir, err)
	}

	// Remove existing linearization files
	linFiles, err := filepath.Glob(filepath.Join(caseDir, "*.lin"))
	if err != nil {
		return nil, err
	}
	for _, linFile := range linFiles {
		os.Remove(linFile)
	}

	// Wrap app context with cancel function
	ctx, cancelFunc := context.WithCancelCause(a.ctx)

	// Save cancel function so it can be called
	EvalCancel = cancelFunc

	// Wrap cancel context with error group so eval will stop on first error
	g, ctx2 := errgroup.WithContext(ctx)

	// Create eval status slice
	statuses := []EvalStatus{}

	// Launch evaluations throttled to number of CPUs specified
	semChan := make(chan struct{}, numCPUs)
	for _, op := range c.OperatingPoints {
		op := op
		statuses = append(statuses, EvalStatus{ID: op.ID, State: "Queued"})
		g.Go(func() error {
			semChan <- struct{}{}
			defer func() { <-semChan }()
			return a.Project.EvaluateLinearization(ctx2, c, &op, caseDir)
		})
	}

	// Wait for evaluations to complete. If error, print
	go func() {
		if err := g.Wait(); err != nil {
			runtime.LogErrorf(a.ctx, "error evaluating case: %s", err)
		}
		cancelFunc(nil) // cancel the context
	}()

	return statuses, nil
}

func (a *App) CancelEvaluate() {
	EvalCancel(fmt.Errorf("evaluation canceled"))
}

//------------------------------------------------------------------------------
// Results
//------------------------------------------------------------------------------

func (a *App) OpenCaseDirectoryDialog() (*Project, error) {

	casePath := strings.TrimSuffix(a.Project.Info.Path, filepath.Ext(a.Project.Info.Path))

	// Open dialog so user can select the case directory
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            "Open Case Directory",
		DefaultDirectory: casePath,
	})
	if err != nil {
		return nil, fmt.Errorf("error selecting OpenFAST file: %w", err)
	}

	// Search for linearization files
	LinFiles, err := filepath.Glob(filepath.Join(path, "*.lin"))
	if err != nil {
		return nil, err
	}
	if len(LinFiles) == 0 {
		return nil, fmt.Errorf("no linearization files found")
	}

	// Process linearization files into results
	a.Project.Results, err = LoadResults(LinFiles)
	if err != nil {
		return nil, err
	}

	// Save results to file
	bs, err := json.MarshalIndent(a.Project.Results, "", "\t")
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(filepath.Join(path, "results.json"), bs, 0777)
	if err != nil {
		return nil, err
	}

	// Initialize return Project
	p := &Project{Info: a.Project.Info, Results: a.Project.Results}

	// Parse and return model
	return p, nil
}
