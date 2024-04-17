package lin_test

import (
	"acdc/lin"
	"path/filepath"
	"testing"
)

func TestLoadResults(t *testing.T) {

	// dir := "../samples/NREL_5MW-ED/case01"
	// dir := "../samples/autoModeTrackingModels/nrel5mw/structOnly75_fast"
	// dir := "../samples/autoModeTrackingModels/nrel5mw/aeroStructSteady"
	// dir := "../samples/autoModeTrackingModels/iea15mw/structOnly75"
	dir := "testdata/large"

	LinFiles, err := filepath.Glob(filepath.Join(dir, "*.lin"))
	if err != nil {
		t.Fatal(err)
	}
	if len(LinFiles) == 0 {
		t.Fatal("no lin files found")
	}

	_, err = lin.ProcessFiles(LinFiles)
	if err != nil {
		t.Fatal(err)
	}
}
