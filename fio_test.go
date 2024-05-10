package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestFiles(t *testing.T) {

	const fstPath = "testdata/fio/NREL_5MW.fst"

	// Parse main file and associated files
	files, err := ParseFiles(fstPath)
	if err != nil {
		t.Fatalf("error parsing '%s': %s", fstPath, err)
	}

	// Create output directory to write files to
	path := filepath.Join("testdata", "output", "fio")
	if err := os.MkdirAll(path, 0777); err != nil {
		t.Fatal(err)
	}

	// Use files structure to recreate input files
	if err := files.Write(path, ""); err != nil {
		t.Fatal(err)
	}

	// Write file data as JSON file
	b, err := json.MarshalIndent(files, "", "\t")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(path, "files.json"), b, 0777); err != nil {
		t.Fatal(err)
	}

	// Check that it parsed the appropriate files
	if len(files.Main) != 1 {
		t.Fatal("Failed to read fst file")
	}
	if len(files.ElastoDyn) != 1 {
		t.Fatal("Failed to read ElastoDyn file")
	}
	if len(files.BeamDyn) != 1 {
		t.Fatal("Failed to read BeamDyn file")
	}
	if len(files.ServoDyn) != 1 {
		t.Fatal("Failed to read ServoDyn file")
	}
	if len(files.AeroDyn) != 1 {
		t.Fatal("Failed to read AeroDyn file")
	}
	if len(files.Misc) != 4 {
		t.Fatal("Failed to read all Misc files")
	}
}
