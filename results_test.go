package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadResults(t *testing.T) {

	dir := "lin/testdata/bd_aero"

	res, err := ProcessCaseDir(dir)
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
