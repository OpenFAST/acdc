package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestFiles(t *testing.T) {

	inputFileDirs := []string{
		"fio-v3.5.x",
		"fio-v4.0.x",
	}

	for _, dir := range inputFileDirs {

		fstPath := filepath.Join("testdata", dir, "NREL_5MW.fst")

		// Parse main file and associated files
		files, err := ParseFiles(fstPath)
		if err != nil {
			t.Fatalf("error parsing '%s': %s", fstPath, err)
		}

		// Create output directory to write files to
		outPath := filepath.Join("testdata", "output", dir)
		if err := os.MkdirAll(outPath, 0777); err != nil {
			t.Fatal(err)
		}

		// Use files structure to recreate input files
		if err := files.Write(outPath, ""); err != nil {
			t.Fatal(err)
		}

		// Write file data as JSON file
		b, err := json.MarshalIndent(files, "", "\t")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(outPath, "files.json"), b, 0777); err != nil {
			t.Fatal(err)
		}

		// Check that it parsed the appropriate files
		if act, exp := len(files.Main), 1; act != exp {
			t.Fatalf("Expected %d Main files, got %d", exp, act)
		}
		if act, exp := len(files.ElastoDyn), 1; act != exp {
			t.Fatalf("Expected %d ElastoDyn files, got %d", exp, act)
		}
		if act, exp := len(files.BeamDyn), 1; act != exp {
			t.Fatalf("Expected %d BeamDyn files, got %d", exp, act)
		}
		if act, exp := len(files.ServoDyn), 1; act != exp {
			t.Fatalf("Expected %d ServoDyn files, got %d", exp, act)
		}
		if act, exp := len(files.InflowWind), 1; act != exp {
			t.Fatalf("Expected %d InflowWind files, got %d", exp, act)
		}
		if act, exp := len(files.AeroDyn), 1; act != exp {
			t.Fatalf("Expected %d AeroDyn files, got %d", exp, act)
		}
		if act, exp := len(files.Misc), 4; act != exp {
			t.Fatalf("Expected %d Misc files, got %d", exp, act)
		}
	}

}
