package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	goruntime "runtime"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/sync/errgroup"
)

type Evaluate struct {
	ExecPath    string `json:"ExecPath"`
	ExecVersion string `json:"ExecVersion"`
	ExecValid   bool   `json:"ExecValid"`
	MaxCPUs     int    `json:"MaxCPUs"`
	NumCPUs     int    `json:"NumCPUs"`
	FilesOnly   bool   `json:"FilesOnly"`
}

var SendEvalStatus = func(ctx context.Context, es EvalStatus) {
	runtime.EventsEmit(ctx, "evalStatus", es)
}

func NewEvaluate() *Evaluate {
	return &Evaluate{
		NumCPUs: 1,
		MaxCPUs: goruntime.NumCPU(),
	}
}

type EvalStatus struct {
	ID          int    `json:"ID"`
	State       string `json:"State"`
	SimProgress int    `json:"SimProgress"`
	LinProgress int    `json:"LinProgress"`
	LogPath     string `json:"LogPath"`
	Error       string `json:"Error"`
}

var EvalCancel context.CancelCauseFunc = func(_ error) {}

func (eval *Evaluate) Case(appCtx context.Context, model *Model, c *Case, projectRootPath string) ([]EvalStatus, error) {

	// Call existing cancel func
	EvalCancel(fmt.Errorf("new evaluation started"))

	// Create path to case directory
	caseDir := filepath.Join(projectRootPath, fmt.Sprintf("Case%02d", c.ID))
	if err := os.MkdirAll(caseDir, 0777); err != nil {
		return nil, fmt.Errorf("error creating directory '%s': %w", caseDir, err)
	}

	// Remove existing output files
	extsToRemove := map[string]struct{}{".lin": {}, ".stamp": {}, ".out": {}, ".vtp": {}}
	filepath.WalkDir(caseDir, func(path string, d fs.DirEntry, err error) error {
		if _, ok := extsToRemove[filepath.Ext(path)]; ok {
			os.Remove(path)
		}
		return nil
	})

	// Wrap app context with cancel function
	ctx, cancelFunc := context.WithCancelCause(appCtx)

	// Save cancel function so it can be called
	EvalCancel = cancelFunc

	// Wrap cancel context with error group so eval will stop on first error
	g, ctx2 := errgroup.WithContext(ctx)

	// Create eval status slice
	statuses := []EvalStatus{}

	// Launch evaluations throttled to number of CPUs specified
	semChan := make(chan struct{}, eval.NumCPUs)
	for _, op := range c.OperatingPoints {
		op := op
		statuses = append(statuses, EvalStatus{ID: op.ID, State: "Queued"})
		g.Go(func() error {
			<-semChan
			defer func() { semChan <- struct{}{} }()
			return eval.OP(ctx2, model, c, &op, caseDir)
		})
	}

	// Wait for evaluations to complete. If error, print
	go func() {

		// Get error from group
		err := g.Wait()
		if err != nil {
			runtime.LogErrorf(appCtx, "error evaluating case: %s", err)
		}

		// Close semaphore channel
		close(semChan)

		// Drain channel
		for {
			if _, ok := <-semChan; !ok {
				break
			}
		}

		// Cancel the context for cleanup
		cancelFunc(nil)

		// If no error, write timestamp of evaluation completion
		if err == nil {
			os.WriteFile(filepath.Join(caseDir, "complete.stamp"),
				[]byte(time.Now().Format(time.RFC3339)), 0777)
		}
	}()

	// Start evaluations
	go func() {
		time.Sleep(time.Second)
		for i := 0; i < eval.NumCPUs; i++ {
			semChan <- struct{}{}
		}
	}()

	return statuses, nil
}

func (eval *Evaluate) OP(ctx context.Context, model *Model, c *Case, op *Condition, caseDir string) error {

	//--------------------------------------------------------------------------
	// Prepare input files
	//--------------------------------------------------------------------------

	// Create a local copy of the files so modifications don't affect the original
	files, err := model.Files.Copy()
	if err != nil {
		return err
	}

	// Check that main and ElastoDyn files exist, return error if not
	if len(files.Main) == 0 {
		return fmt.Errorf("no Main files were imported")
	}
	if len(files.ElastoDyn) == 0 {
		return fmt.Errorf("no ElastoDyn files were imported")
	}

	// Set status update time so a full simulation will generate 50 status messages
	// which are used for the progress bars on the Evaluate tab
	if statusTime := files.Main[0].TMax.Value / 100; statusTime < files.Main[0].SttsTime.Value {
		files.Main[0].SttsTime.Value = statusTime
	}

	// The file writing code assumes that NLinTimes is governed by the
	// length of LinTimes which is only true if CalcSteady==false.
	// If CalcSteady is true, the size of LinTimes must be changed
	// to match NLinTimes
	if files.Main[0].CalcSteady.Value {
		files.Main[0].LinTimes.Value = make([]float64, files.Main[0].NLinTimes.Value)
	}

	// VTK Options
	files.Main[0].WrVTK.Value = 3    // Write checkpoint files for mode shape generation
	files.Main[0].VTK_type.Value = 2 // Write lines meshes to VTK files
	files.Main[0].VTK_fps.Value = 24 // Generate 24 frames of mode animation

	// Modify ElastoDyn file for operating point conditions
	files.ElastoDyn[0].BlPitch1.Value = op.BladePitch
	files.ElastoDyn[0].BlPitch2.Value = op.BladePitch
	files.ElastoDyn[0].BlPitch3.Value = op.BladePitch
	files.ElastoDyn[0].RotSpeed.Value = op.RotorSpeed

	// Set RotateStates to true for all BeamDyn files
	for i := range files.BeamDyn {
		files.BeamDyn[i].RotStates.Value = true
	}

	// If case includes aero and wind speed is nonzero
	if c.IncludeAero && op.WindSpeed > 0 {

		// Set flag to use InflowWind
		files.Main[0].CompInflow.Value = 1

		// Set flag to use AeroDyn or return error
		if len(files.AeroDyn) > 0 {
			files.Main[0].CompAero.Value = 2
		} else if len(files.AeroDyn14) > 0 {
			files.Main[0].CompAero.Value = 1
		} else {
			return fmt.Errorf("no Aero files were imported")
		}

		// Set flag to use InflowWind or return error
		if len(files.InflowWind) == 0 {
			return fmt.Errorf("no InflowWind files were imported")
		}
		files.InflowWind[0].WindType.Value = 1
		files.InflowWind[0].HWindSpeed.Value = op.WindSpeed
		files.InflowWind[0].PLExp.Value = 0

		// If controller is enabled
		if c.UseController {

			// If no ServoDyn file, return error
			if len(files.ServoDyn) == 0 {
				return fmt.Errorf("no ServoDyn files were imported")
			}

			// Set CompServo to 1
			files.Main[0].CompServo.Value = 1

			// Enable generator DOF
			files.ElastoDyn[0].GenDOF.Value = true

			// Set ServoDyn parameters
			files.ServoDyn[0].PCMode.Value = 0
			files.ServoDyn[0].VSContrl.Value = 1
			files.ServoDyn[0].HSSBrMode.Value = 0
			files.ServoDyn[0].YCMode.Value = 0

			// Set trim mode based on below or above rated wind speed
			if op.WindSpeed < float64(c.RatedWindSpeed) {
				files.Main[0].TrimCase.Value = 2
				files.Main[0].TrimGain.Value = c.TrimGain[0]
				files.ServoDyn[0].VS_RtGnSp.Value = c.RatedRotorSpeed
			} else {
				files.Main[0].TrimCase.Value = 3
				files.Main[0].TrimGain.Value = c.TrimGain[1]
				files.ServoDyn[0].VS_RtGnSp.Value = 1e-3
			}

		} else {

			// Disable ServoDyn and remove files
			files.Main[0].CompServo.Value = 0
			files.ServoDyn = []ServoDyn{}

			// Disable generator DOF
			files.ElastoDyn[0].GenDOF.Value = false
		}

	} else {

		// Disable InflowWind and remove files
		files.Main[0].CompInflow.Value = 0
		files.InflowWind = []InflowWind{}

		// Disable AeroDyn and remove files
		files.Main[0].CompAero.Value = 0
		files.AeroDyn = []AeroDyn{}
		files.AeroDyn14 = []AeroDyn14{}

		// Disable ElastoDyn generator DOF
		files.ElastoDyn[0].GenDOF.Value = false

		// Disable ServoDyn and remove files
		files.Main[0].CompServo.Value = 0
		files.ServoDyn = []ServoDyn{}
	}

	// Write modified turbine files
	filePrefix := fmt.Sprintf("%02d_", op.ID)
	if err := files.Write(caseDir, filePrefix); err != nil {
		return fmt.Errorf("error writing turbine files: %w", err)
	}

	// Get number of linearization steps
	numLinSteps := files.Main[0].NLinTimes.Value

	// Create path to main file and log file
	mainName := filePrefix + files.Main[0].Name
	rootName := strings.TrimSuffix(mainName, filepath.Ext(mainName))
	rootPath := filepath.Join(caseDir, rootName)
	mainPath := filepath.Join(caseDir, mainName)
	logPath := rootPath + ".log"

	// Create status ID from operating point and set in linearization flag to false
	statusID := op.ID
	inLinearization := false

	// If flag set to only output the files (not run simulation), return
	if eval.FilesOnly {
		return nil
	}

	//--------------------------------------------------------------------------
	// Run Linearization
	//--------------------------------------------------------------------------

	// Create log file
	logFile, err := os.Create(logPath)
	if err != nil {
		return fmt.Errorf("error creating log file '%s': %w", logFile.Name(), err)
	}
	defer logFile.Close()

	// Create command, get output pipe, set stderr to stdout, start command
	cmd := exec.CommandContext(ctx, eval.ExecPath, mainPath)
	outputReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = cmd.Stdout
	cmd.Start()

	// Get progress
	scanner := bufio.NewScanner(outputReader)
	scanner.Split(ScanLines)
	for scanner.Scan() {
		line := strings.Map(func(r rune) rune {
			if unicode.IsGraphic(r) {
				return r
			}
			return -1
		}, scanner.Text())
		line = strings.TrimSpace(line)
		logFile.WriteString(line + "\n")
		if strings.Contains(line, "Time: ") && !inLinearization {
			fields := strings.Fields(line)
			currentTime, err := strconv.ParseFloat(fields[1], 32)
			if err != nil {
				continue
			}
			totalTime, err := strconv.ParseFloat(fields[3], 32)
			if err != nil {
				continue
			}
			SendEvalStatus(ctx, EvalStatus{
				ID:          statusID,
				State:       "Simulation",
				SimProgress: int(100 * currentTime / totalTime),
				LogPath:     logPath,
			})
		} else if strings.Contains(line, "Performing linearization") {
			inLinearization = true
			fields := strings.Fields(line)
			linNumber, err := strconv.ParseFloat(fields[2], 32)
			if err != nil {
				continue
			}
			SendEvalStatus(ctx, EvalStatus{
				ID:          statusID,
				State:       "Linearization",
				SimProgress: 100,
				LinProgress: int(100 * linNumber / float64(numLinSteps)),
				LogPath:     logPath,
			})
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	// Wait for command to exit, set status on error
	if err := cmd.Wait(); err != nil {
		status := EvalStatus{
			ID:          statusID,
			State:       "Error",
			SimProgress: 100,
			LinProgress: 100,
			Error:       err.Error(),
			LogPath:     logPath,
		}
		// If context was canceled, set state to canceled
		if cause := context.Cause(ctx); cause != nil {
			status.State = "Canceled"
			status.Error = cause.Error()
		}
		SendEvalStatus(ctx, status)
		return err
	}

	//--------------------------------------------------------------------------
	// Send complete status
	//--------------------------------------------------------------------------

	SendEvalStatus(ctx, EvalStatus{
		ID:          statusID,
		State:       "Complete",
		SimProgress: 100,
		LinProgress: 100,
		LogPath:     logPath,
	})

	return nil
}

func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexAny(data, "\n\r"); i >= 0 {
		if data[i] == '\n' { // LF line
			return i + 1, data[0:i], nil
		} else if len(data) > i+1 && data[i+1] == '\n' { // CRLF line
			return i + 2, data[0:i], nil
		} else { // CR Line
			return i + 1, data[0:i], nil
		}
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}
