package lin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

type FileGroup struct {
	Name  string
	Files []string
}

type LinOP struct {
	RootPath      string        `json:"Name"`
	FilePaths     []string      `json:"-"`
	HasAeroStates bool          `json:"HasAeroStates"`
	MBC           *MBC          `json:"MBC"`
	EigRes        *EigenResults `json:"EigRes"`
}

// ProcessFiles takes a slice of linearization file paths, groups them by
// operating point and performs MBC/Eigenanalysis for each OP. It returns
// a slice of operating point linearization results.
func ProcessFiles(LinFilePaths []string) ([]LinOP, error) {

	// Group linearization files by operating point using the file name
	opFilesMap := map[string]*FileGroup{}
	for _, filePath := range LinFilePaths {

		// Split file path by '.'
		tmp := strings.Split(filePath, ".")

		// Get the root name of the file without extension or lin numbers
		rootName := strings.Join(tmp[:len(tmp)-2], ".")

		// Get file group for this operating point, create it if it doesn't exist
		opFiles, ok := opFilesMap[rootName]
		if !ok {
			opFiles = &FileGroup{Name: rootName}
			opFilesMap[rootName] = opFiles
		}

		// Append current file path to slice of files
		opFiles.Files = append(opFiles.Files, filePath)
	}

	// Initialize flag that wind was present to false
	hasWind := false

	results := []LinOP{}
	for _, fg := range opFilesMap {
		// Read all linearization files in group
		linFileData := make([]*LinData, len(fg.Files))
		var err error
		for i, linFilePath := range fg.Files {
			linFileData[i], err = ReadLinFile(linFilePath)
			if err != nil {
				return nil, fmt.Errorf("error reading '%s': %w", linFilePath, err)
			}
		}

		// Extract matrix data from linearization file data
		matData := NewMatData(linFileData)

		// Perform multi-blade coordinate transform
		mbc, err := matData.MBC3()
		if err != nil {
			return nil, err
		}

		// Perform Eigenanalysis to get modes
		eigRes, err := mbc.EigenAnalysis()
		if err != nil {
			return nil, err
		}

		// Write MBC data to file
		bs, err := json.MarshalIndent(mbc, "", "\t")
		if err != nil {
			return nil, err
		}
		err = os.WriteFile(fg.Name+"_mbc.json", bs, 0777)
		if err != nil {
			return nil, err
		}

		// Write Eigen analysis mode results data to file
		w := &bytes.Buffer{}
		eigRes.Modes.ToCSV(w)
		err = os.WriteFile(fg.Name+"_modes.csv", w.Bytes(), 0777)
		if err != nil {
			return nil, err
		}

		// Determine if any AeroDyn states are in eigenanalysis
		hasAeroStates := false
		for _, dof := range mbc.DOFsEigen {
			if strings.HasPrefix(dof, "AD") {
				hasAeroStates = true
				break
			}
		}

		// If wind speed is nonzero, set flag to true
		if mbc.WindSpeed > 0 {
			hasWind = true
		}

		// Return MBC and eigen analysis results
		results = append(results, LinOP{
			RootPath:      fg.Name,
			FilePaths:     fg.Files,
			HasAeroStates: hasAeroStates,
			MBC:           mbc,
			EigRes:        eigRes,
		})
	}

	// Sort results by wind speed or rotor speed
	if hasWind {
		sort.Slice(results, func(i, j int) bool {
			return results[i].MBC.WindSpeed < results[j].MBC.WindSpeed
		})
	} else {
		sort.Slice(results, func(i, j int) bool {
			return results[i].MBC.RotSpeed < results[j].MBC.RotSpeed
		})
	}

	// Add operating point IDs to modes now that they've been sorted
	for i := range results {
		for j := range results[i].EigRes.Modes {
			results[i].EigRes.Modes[j].OP = i
		}
	}

	return results, nil
}
