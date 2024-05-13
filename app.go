package main

import (
	"acdc/diagram"
	"acdc/viz"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

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

func (a *App) OpenProjectDialog() (*Info, error) {

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

	// Load project
	a.Project, err = LoadProject(path)
	if err != nil {
		return nil, err
	}

	// Set project path
	a.Project.Info.Path = path

	// set window title
	runtime.WindowSetTitle(a.ctx, "ACDC - "+path)

	// Open project
	return &a.Project.Info, err
}

// SaveProjectDialog opens a dialog to select where to save the project
func (a *App) SaveProjectDialog() (*Info, error) {

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

	// Set project path
	a.Project.Info.Path = path

	// Save project
	if _, err := a.Project.Save(); err != nil {
		return nil, err
	}

	// Return full project
	return &a.Project.Info, nil
}

func (a *App) OpenProject(path string) (*Info, error) {

	// Load project
	var err error
	a.Project, err = LoadProject(path)
	if err != nil {
		return nil, err
	}

	// Set project path
	a.Project.Info.Path = path

	return &a.Project.Info, nil
}

//------------------------------------------------------------------------------
// Model
//------------------------------------------------------------------------------

// UpdateModel saves changes to the current project
func (a *App) FetchModel() (*Model, error) {

	// If no Model in project, create it
	if a.Project.Model == nil {
		a.Project.Model = NewModel()
	}

	return a.Project.Model, nil
}

// UpdateModel saves changes to the current project
func (a *App) UpdateModel(model *Model) (*Model, error) {

	// Update model in project
	a.Project.Model = model

	// Save project
	if _, err := a.Project.Save(); err != nil {
		return nil, err
	}

	return a.Project.Model, nil
}

func (a *App) ImportModelDialog() (*Model, error) {

	// Open dialog so user can select the file
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Open Model",
		Filters: []runtime.FileFilter{
			{DisplayName: "OpenFAST Model (*.fst)", Pattern: "*.fst"},
		},
	})
	if err != nil {
		// No file was selected, return current model
		return a.Project.Model, nil
	}

	// Parse model files
	a.Project.Model, err = ParseModelFiles(path)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error parsing model files: %s", err)
		return nil, err
	}

	// Save project
	if _, err := a.Project.Save(); err != nil {
		runtime.LogErrorf(a.ctx, "SelectExec: error saving project: %s", err)
	}

	// Parse and return model
	return a.Project.Model, nil
}

//------------------------------------------------------------------------------
// Analysis
//------------------------------------------------------------------------------

func (a *App) FetchAnalysis() (*Analysis, error) {

	// If no analysis in project, create it
	if a.Project.Analysis == nil {
		a.Project.Analysis = NewAnalysis()
	}

	// Save project
	if _, err := a.Project.Save(); err != nil {
		return nil, err
	}

	return a.Project.Analysis, nil
}

// UpdateEval runs Case.Calculate and saves evaluation data
func (a *App) UpdateAnalysis(an *Analysis) (*Analysis, error) {

	// Calculate analysis cases
	if err := an.Calculate(); err != nil {
		return nil, err
	}

	// Update analysis in the project
	a.Project.Analysis = an

	// Save project
	if _, err := a.Project.Save(); err != nil {
		return nil, err
	}

	return a.Project.Analysis, nil
}

func (a *App) AddAnalysisCase() (*Analysis, error) {

	// If no analysis in project, create it
	if a.Project.Analysis == nil {
		a.Project.Analysis = NewAnalysis()
	}

	// Add new case to analysis
	a.Project.Analysis.Cases = append(a.Project.Analysis.Cases, NewCase())

	// Calculate analysis cases
	if err := a.Project.Analysis.Calculate(); err != nil {
		return nil, err
	}

	// Save project
	if _, err := a.Project.Save(); err != nil {
		return nil, err
	}

	return a.Project.Analysis, nil
}

func (a *App) DuplicateAnalysisCase(caseID int) (*Analysis, error) {

	// If no analysis in project, create it
	if a.Project.Analysis == nil {
		a.Project.Analysis = NewAnalysis()
	}

	// Find case to copy
	srcCase, err := a.Project.Case(caseID)
	if err != nil {
		return nil, err
	}

	// Copy case
	dstCase, err := srcCase.Copy()
	if err != nil {
		return nil, err
	}

	// Add "copy" to destination case name
	dstCase.Name += " Copy"

	// Append case to analysis
	a.Project.Analysis.Cases = append(a.Project.Analysis.Cases, dstCase)

	// Calculate analysis cases
	if err := a.Project.Analysis.Calculate(); err != nil {
		return nil, err
	}

	// Save project
	if _, err := a.Project.Save(); err != nil {
		return nil, err
	}

	return a.Project.Analysis, nil
}

func (a *App) RemoveAnalysisCase(caseID int) (*Analysis, error) {

	// If no analysis in project, create it
	if a.Project.Analysis == nil {
		a.Project.Analysis = NewAnalysis()
	}

	// Filter out case that matches ID
	tmp := []Case{}
	for _, c := range a.Project.Analysis.Cases {
		if c.ID != caseID {
			tmp = append(tmp, c)
		}
	}
	a.Project.Analysis.Cases = tmp

	// Calculate analysis cases
	if err := a.Project.Analysis.Calculate(); err != nil {
		return nil, err
	}

	// Save project
	if _, err := a.Project.Save(); err != nil {
		return nil, err
	}

	return a.Project.Analysis, nil
}

func (a *App) ImportAnalysisCaseCurve(caseID int) (*Analysis, error) {

	// If no analysis in project, create it
	if a.Project.Analysis == nil {
		a.Project.Analysis = NewAnalysis()
	}

	// Allow use to select the file
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select curve data file",
		Filters: []runtime.FileFilter{
			{DisplayName: "CSV (*.csv)", Pattern: "*.csv"},
		},
	})
	if err != nil {
		return a.Project.Analysis, nil
	}

	// Open file
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening '%s': %w", path, err)
	}
	defer f.Close()

	// Create CSV reader and read data
	cr := csv.NewReader(f)
	cr.Comment = '#'
	rows, err := cr.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error parsing '%s': %w", path, err)
	}

	// If no rows read, return error
	if len(rows) == 0 {
		return nil, fmt.Errorf("data file '%s': is empty", path)
	}

	// Check that there are at least 3 columns
	if n := len(rows[0]); n < 3 {
		return nil, fmt.Errorf("data file '%s' has %d columns, 3 are required", path, n)
	}

	// Get pointer to the case or return error if invalid case ID
	if caseID < 1 || caseID > len(a.Project.Analysis.Cases) {
		return nil, fmt.Errorf("invalid case ID: %d", caseID)
	}
	c := &a.Project.Analysis.Cases[caseID-1]

	// Initialize curve to length of rows
	c.Curve = make([]Condition, len(rows))

	// Loop through rows and populate curve
	for i, row := range rows {

		c.Curve[i].WindSpeed, err = strconv.ParseFloat(strings.TrimSpace(row[0]), 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing wind speed from '%s': %w", row[0], err)
		}
		c.Curve[i].RotorSpeed, err = strconv.ParseFloat(strings.TrimSpace(row[1]), 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing rotor speed from '%s': %w", row[1], err)
		}
		c.Curve[i].BladePitch, err = strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing blade pitch from '%s': %w", row[2], err)
		}
	}

	// Calculate analysis cases
	if err := a.Project.Analysis.Calculate(); err != nil {
		return nil, err
	}

	// Save project
	if _, err := a.Project.Save(); err != nil {
		return nil, err
	}

	return a.Project.Analysis, nil
}

//------------------------------------------------------------------------------
// Evaluate
//------------------------------------------------------------------------------

func (a *App) FetchEvaluate() (*Evaluate, error) {

	newEvaluate := NewEvaluate()

	// If no Eval in project
	if a.Project.Evaluate == nil {

		// Set project value from new evaluate
		a.Project.Evaluate = newEvaluate
	} else {

		// Get max CPUs from new evaluate
		a.Project.Evaluate.MaxCPUs = newEvaluate.MaxCPUs
	}

	// Save project
	if _, err := a.Project.Save(); err != nil {
		return nil, err
	}

	return a.Project.Evaluate, nil
}

// UpdateEvaluate runs Case.Calculate and saves evaluation data
func (a *App) UpdateEvaluate(ev *Evaluate) (*Evaluate, error) {

	// Update analysis in the project
	a.Project.Evaluate = ev

	// Save project
	if _, err := a.Project.Save(); err != nil {
		return nil, err
	}

	return a.Project.Evaluate, nil
}

// SelectExec opens a dialog for the user to select an OpenFAST executable.
func (a *App) SelectExec() (*Evaluate, error) {

	// If no Eval in project, create it
	if a.Project.Evaluate == nil {
		a.Project.Evaluate = NewEvaluate()
	}

	// Lookup OpenFAST executable with default name
	execPath, err := exec.LookPath("openfast")
	if err != nil {
		execPath = ""
	}

	// Open dialog so user can select the executable path
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:                "Select OpenFAST Executable",
		CanCreateDirectories: false,
		DefaultDirectory:     filepath.Dir(execPath),
	})
	if err != nil {
		return nil, err
	}

	// If no path selected, return current exec
	if path == "" {
		return a.Project.Evaluate, nil
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
	a.Project.Evaluate.ExecPath = path
	a.Project.Evaluate.ExecVersion = string(bytes.TrimSpace(output))
	a.Project.Evaluate.ExecValid = true

	// Save project
	if _, err := a.Project.Save(); err != nil {
		runtime.LogErrorf(a.ctx, "SelectExec: error saving project: %s", err)
	}

	// Save project and return path
	return a.Project.Evaluate, nil
}

func (a *App) EvaluateCase(caseID int) ([]EvalStatus, error) {

	// Find case
	c, err := a.Project.Case(caseID)
	if err != nil {
		return []EvalStatus{}, err
	}

	// Evaluate case
	return a.Project.Evaluate.Case(a.ctx, a.Project.Model, c, a.Project.RootPath())
}

func (a *App) CancelEvaluate() {
	EvalCancel(fmt.Errorf("evaluation canceled"))
}

func (a *App) GetEvaluateLog(path string) (string, error) {
	bs, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

//------------------------------------------------------------------------------
// Results
//------------------------------------------------------------------------------

type LinDirData struct {
	Dir     string           `json:"Dir"`
	Results *Results         `json:"Results"`
	Diagram *diagram.Diagram `json:"Diagram"`
}

func (a *App) SelectCaseLinDir(caseID int) (LinDirData, error) {

	// Find case
	c, err := a.Project.Case(caseID)
	if err != nil {
		return LinDirData{}, err
	}

	// Build case directory
	linDir := filepath.Join(a.Project.RootPath(), fmt.Sprintf("Case%02d", c.ID))

	// Create linearization directory data structure
	ld := LinDirData{Dir: linDir}

	// Attempt to load results and diagram
	if results, err := LoadResults(linDir); err == nil {
		a.Project.Results = results
		ld.Results = results.ForApp()
	} else {
		fmt.Println(err)
	}
	if diag, err := diagram.Load(filepath.Join(linDir, "diagram.json")); err == nil {
		a.Project.Diagram = diag
		ld.Diagram = diag
	} else {
		fmt.Println(err)
	}

	return ld, nil
}

func (a *App) SelectCustomLinDir() (LinDirData, error) {

	// Get path to project, if it doesn't exist, set to empty string
	projectDir := filepath.Dir(a.Project.Info.Path)
	if _, err := os.Stat(projectDir); err != nil {
		projectDir = ""
	}

	// Open dialog so user can select the case directory
	linDir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            "Open Case Directory",
		DefaultDirectory: projectDir,
	})
	if err != nil && linDir != "" {
		return LinDirData{}, err
	}

	// Create linearization directory data structure
	ld := LinDirData{Dir: linDir}

	// Attempt to load results and diagram
	if results, err := LoadResults(linDir); err == nil {
		a.Project.Results = results
		ld.Results = results.ForApp()
	} else {
		fmt.Println(err)
	}
	if diag, err := diagram.Load(filepath.Join(linDir, "diagram.json")); err == nil {
		a.Project.Diagram = diag
		ld.Diagram = diag
	} else {
		fmt.Println(err)
	}

	return ld, nil
}

func (a *App) FetchResults() (*Results, error) {

	// If no Results in project, create it
	if a.Project.Results == nil {
		a.Project.Results = NewResults()
	}

	// Save project
	if _, err := a.Project.Save(); err != nil {
		return nil, err
	}

	return a.Project.Results.ForApp(), nil
}

func (a *App) ProcessLinDir(linDir string) (*Results, error) {

	// Process case directory to get results
	results, err := ProcessCaseDir(linDir)
	if err != nil {
		return nil, err
	}

	// Add results to project
	a.Project.Results = results

	// Save project
	if _, err := a.Project.Save(); err != nil {
		return nil, err
	}

	// Save results
	if err := results.Save(linDir); err != nil {
		return nil, err
	}

	return a.Project.Results.ForApp(), nil
}

//------------------------------------------------------------------------------
// Diagram
//------------------------------------------------------------------------------

func (a *App) GenerateDiagram(opts diagram.Options) (*diagram.Diagram, error) {

	// Check that results have been loaded
	if a.Project.Results == nil {
		return nil, fmt.Errorf("load results before generating diagram")
	}

	// Generate diagram with given options
	diag, err := diagram.New(a.Project.Results.LinOPs, opts)
	if err != nil {
		return nil, err
	}

	// Save diagram in project
	a.Project.Diagram = diag

	// Save project
	if _, err := a.Project.Save(); err != nil {
		return nil, err
	}

	return a.Project.Diagram, nil
}

// UpdateDiagram saves the diagram to file
func (a *App) UpdateDiagram(diag *diagram.Diagram) error {

	// Update analysis in the project
	a.Project.Diagram = diag

	// Save project
	if _, err := a.Project.Save(); err != nil {
		return err
	}

	return nil
}

//------------------------------------------------------------------------------
// Visualization
//------------------------------------------------------------------------------

func (a *App) GetModeViz(opID int, modeID int, scale float32) (*viz.ModeData, error) {

	// If executable path is not valid, return error
	if _, err := exec.LookPath(a.Project.Evaluate.ExecPath); err != nil {
		return nil, fmt.Errorf("executable path is not valid")
	}

	// If results haven't been loaded, return error
	if a.Project.Results == nil {
		return nil, fmt.Errorf("load results before generating visualization")
	}

	// If operating point index is not valid, return error
	if opID < 0 || opID >= len(a.Project.Results.LinOPs) {
		return nil, fmt.Errorf("invalid operating point ID: %d", opID)
	}

	// If mode index is not valid, return error
	if modeID < 0 || modeID >= len(a.Project.Results.LinOPs[opID].Modes) {
		return nil, fmt.Errorf("invalid mode ID (%d) for operating point (%d)", modeID, opID)
	}

	// Create visualization options
	opts := viz.Options{Scale: scale}

	// Generate mode visualization data
	modeData, err := opts.GenerateModeData(a.Project.Evaluate.ExecPath,
		&a.Project.Results.LinOPs[opID], []int{modeID})
	if err != nil {
		return nil, err
	}

	// Populate operating point ID and mode ID
	modeData.OPID = opID
	modeData.ModeID = modeID

	return modeData, nil
}

//------------------------------------------------------------------------------
// Export Data
//------------------------------------------------------------------------------

func (a *App) ExportDiagramDataJSON(diag diagram.Diagram) error {

	// Open dialog so user can select the file
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:                "Save Diagram Data As",
		DefaultFilename:      "campbell_diagram.json",
		CanCreateDirectories: true,
	})
	if err != nil {
		return nil
	}

	// Create path
	if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		return fmt.Errorf("error creating project directory '%s': %w", path, err)
	}

	// Convert config into JSON
	bs, err := json.MarshalIndent(diag, "", "\t")
	if err != nil {
		return fmt.Errorf("error marshalling data: %w", err)
	}

	// Write file and return error
	return os.WriteFile(path, bs, 0777)
}
