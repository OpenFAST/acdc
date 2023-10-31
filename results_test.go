package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadResults(t *testing.T) {

	// dir := "../samples/NREL_5MW-ED/case01"
	// dir := "../samples/autoModeTrackingModels/nrel5mw/structOnly75_fast"
	// dir := "../samples/autoModeTrackingModels/nrel5mw/aeroStructSteady"
	dir := "../samples/autoModeTrackingModels/iea15mw/structOnly75"

	LinFiles, err := filepath.Glob(filepath.Join(dir, "*.lin"))
	if err != nil {
		t.Fatal(err)
	}
	if len(LinFiles) == 0 {
		t.Fatal("no lin files found")
	}

	res := Results{}
	if err = res.ProcessFiles(LinFiles); err != nil {
		t.Fatal(err)
	}

	// Save results to file
	bs, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(dir, "results.json"), bs, 0777)
	if err != nil {
		t.Fatal(err)
	}
}
