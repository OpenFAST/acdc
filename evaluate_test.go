package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestEvaluate(t *testing.T) {

	SendEvalStatus = func(ctx context.Context, es EvalStatus) {}

	// Load example project
	project, err := LoadProject("testdata/eval/NREL-5MW.json")
	if err != nil {
		t.Fatal(err)
	}

	// Get path to OpenFAST executable, if not found skip test
	execPath, found := os.LookupEnv("OPENFAST_PATH")
	if !found {
		t.Skip("executable path not specified in OPENFAST_PATH environment variable")
	}

	// Get first case
	c := project.Analysis.Cases[0]

	// Get operating point, set rotor speed to zero so linearization happens immediately
	op := c.OperatingPoints[0]
	op.RotorSpeed = 0

	// Create path to case directory
	caseDir := filepath.Join("testdata/output/eval", fmt.Sprintf("Case%02d", c.ID))

	// Remove case directory if it exists
	err = os.RemoveAll(caseDir)
	if err != nil {
		t.Fatal(err)
	}

	// Create case directory
	if err := os.MkdirAll(caseDir, 0777); err != nil {
		t.Fatal(err)
	}

	// Create evaluate struct
	eval := &Evaluate{
		ExecPath:  execPath,
		ExecValid: true,
		FilesOnly: false,
	}

	// Wrap app context with cancel function
	ctx := context.Background()

	// Evaluate operating point
	err = eval.OP(ctx, project.Model, &c, &op, caseDir)
	if err != nil {
		t.Fatal(err)
	}
}
