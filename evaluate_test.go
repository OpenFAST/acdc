package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEvaluate(t *testing.T) {

	project, err := LoadProject("/Users/dslaught/Downloads/PPI_aero.json")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(project)

	c := &project.Analysis.Cases[0]

	// Create path to case directory
	caseDir := filepath.Join(strings.TrimSuffix(project.Info.Path, filepath.Ext(project.Info.Path)), fmt.Sprintf("case%02d", c.ID))
	if err := os.MkdirAll(caseDir, 0777); err != nil {
		t.Fatal(err)
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
	ctx := context.Background()

	err = EvaluateOP(ctx, project.Model, c, &c.OperatingPoints[0], caseDir,
		project.Evaluate.ExecPath)
	if err != nil {
		t.Fatal(err)
	}
}
