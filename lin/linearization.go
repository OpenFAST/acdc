package lin

import (
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

type OPResult struct {
	Name          string        `json:"Name"`
	Files         []string      `json:"-"`
	HasAeroStates bool          `json:"HasAeroStates"`
	MBC           *MBC          `json:"MBC"`
	EigRes        *EigenResults `json:"EigRes"`
}

func ProcessFiles(LinFiles []string) ([]OPResult, error) {

	// Group linearization files by operating point using the file name
	fileGroupMap := map[string]*FileGroup{}
	for _, filePath := range LinFiles {
		tmp := strings.Split(filePath, ".")
		filePathNoLinExt := strings.Join(tmp[:len(tmp)-2], ".")
		fileGroup, ok := fileGroupMap[filePathNoLinExt]
		if !ok {
			fileGroup = &FileGroup{Name: filePathNoLinExt}
		}
		fileGroup.Files = append(fileGroup.Files, filePath)
		fileGroupMap[filePathNoLinExt] = fileGroup
	}

	// Initialize flag that wind was present to false
	hasWind := false

	results := []OPResult{}
	for _, fg := range fileGroupMap {
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

		// Write Eigen analysis results data to file
		bs, err = json.MarshalIndent(eigRes.Modes, "", "\t")
		if err != nil {
			return nil, err
		}
		err = os.WriteFile(fg.Name+"_modes.json", bs, 0777)
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
		results = append(results, OPResult{
			Name:          fg.Name,
			Files:         fg.Files,
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
