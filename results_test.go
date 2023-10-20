package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadResults(t *testing.T) {

	// dir := "samples/NREL_5MW-ED/case01"
	dir := "samples/structOnly75_fast"

	LinFiles, err := filepath.Glob(filepath.Join(dir, "*.lin"))
	if err != nil {
		t.Fatal(err)
	}

	res, err := LoadResults(LinFiles)
	if err != nil {
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
