package main

import (
	"fmt"
	"path/filepath"
	"sort"
)

type Model struct {
	HasAero       bool     `json:"HasAero"`
	ImportedPaths []string `json:"ImportedPaths"`
	Files         *Files   `json:"Files"`
	Notes         []string `json:"Notes"`
}

func NewModel() *Model {
	return &Model{}
}

func ParseModelFiles(path string) (*Model, error) {

	// Parse model files
	files, err := ParseFiles(path)
	if err != nil {
		return nil, fmt.Errorf("error importing OpenFAST model '%s': %w", path, err)
	}

	// Get imported paths
	paths := []string{}
	for p := range files.PathMap {
		p, _ := filepath.Rel(filepath.Dir(path), p)
		paths = append(paths, p)
	}
	sort.Slice(paths, func(i, j int) bool {
		return filepath.Base(paths[i]) < filepath.Base(paths[j])
	})

	// Add notes from parsing
	notes := []string{}
	if len(files.AeroDyn)+len(files.AeroDyn14) == 0 {
		notes = append(notes, "No AeroDyn or AeroDyn 14 files imported: aerodynamics option will be disabled in cases")
	}
	if len(files.InflowWind) == 0 {
		notes = append(notes, "No InflowWind file imported: aerodynamics option will be disabled in cases")
	}

	// Initialize models structure
	model := Model{
		HasAero: ((len(files.AeroDyn)+len(files.AeroDyn14)) > 0 &&
			len(files.InflowWind) > 0),
		Files:         files,
		ImportedPaths: paths,
		Notes:         notes,
	}

	return &model, nil
}
