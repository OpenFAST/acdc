package main

import (
	"testing"
)

func TestFiles(t *testing.T) {

	// // Find all FAST files in test-data directory
	// matches, err := filepath.Glob("testdata/fio/*.fst")
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// // Loop through .fst files and parse them
	// for _, fstPath := range matches {

	// 	// Parse main file and associated files
	// 	files, err := ParseFiles(fstPath)
	// 	if err != nil {
	// 		t.Fatalf("error parsing '%s': %s", fstPath, err)
	// 	}

	// 	// Get model name
	// 	modelName := filepath.Base(fstPath)
	// 	modelName = strings.TrimSuffix(modelName, filepath.Ext(modelName))

	// 	// Create output directory to write files to
	// 	path := filepath.Join("testdata", "output", modelName)
	// 	if err := os.MkdirAll(path, 0777); err != nil {
	// 		t.Fatal(err)
	// 	}

	// 	// Use files structure to recreate input files
	// 	if err := files.Write(path, ""); err != nil {
	// 		t.Fatal(err)
	// 	}

	// 	// Write file data as JSON file
	// 	b, err := json.MarshalIndent(files, "", "\t")
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	if err := os.WriteFile(filepath.Join(path, "files.json"), b, 0777); err != nil {
	// 		t.Fatal(err)
	// 	}
	// }

	// Parse main file and associated files
	files, err := ParseFiles("../WindFloatStab/OpenFAST/IEA-15-240-RWT-WindFloat/IEA-15-240-RWT-WindFloat.fst")
	if err != nil {
		t.Fatalf("error parsing '%s': %s", "fstPath", err)
	}

	t.Log(files)
}
