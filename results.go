package main

import (
	"acdc/lin"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type Results struct {
	LinDir  string           `json:"LinDir"`
	HasWind bool             `json:"HasWind"`
	OPs     []OperatingPoint `json:"OPs"`
	LinOPs  []lin.LinOP      `json:"LinOPs"`
}

type OperatingPoint struct {
	ID        int      `json:"ID"`
	Files     []string `json:"Files"`
	RotSpeed  float32  `json:"RotSpeed"`  // RPM
	WindSpeed float32  `json:"WindSpeed"` // m/s
	Modes     []Mode   `json:"Modes"`
}

type Mode struct {
	ID            int     `json:"ID"`
	OP            int     `json:"OP"`
	NaturalFreqHz float32 `json:"NaturalFreqHz"`
	DampedFreqHz  float32 `json:"DampedFreqHz"`
	DampingRatio  float32 `json:"DampingRatio"`
}

func NewResults() *Results {
	return &Results{}
}

// ForApp returns the results structure without the LinOPs member.
func (res *Results) ForApp() *Results {
	results := *res
	results.LinOPs = nil
	return &results
}

func ProcessCaseDir(path string) (*Results, error) {

	// Search for linearization files
	LinFiles, err := filepath.Glob(filepath.Join(path, "*.lin"))
	if err != nil {
		return nil, err
	}
	linRe := regexp.MustCompile(`.+?\.\d+\.lin`)
	tmp := LinFiles
	LinFiles = []string{}
	for _, f := range tmp {
		if linRe.MatchString(f) {
			LinFiles = append(LinFiles, f)
		}
	}
	if len(LinFiles) == 0 {
		return nil, fmt.Errorf("no linearization files found")
	}

	// Process linearization files into results
	linResults, err := lin.ProcessFiles(LinFiles)
	if err != nil {
		return nil, err
	}

	// Initialize results structure
	results := &Results{
		LinDir: path,
		LinOPs: linResults,
	}

	// Extract data from linearization results
	for i, lr := range linResults {
		results.HasWind = lr.MBC.WindSpeed > 0 || results.HasWind
		modes := []Mode{}
		for j, m := range lr.Modes {
			modes = append(modes, Mode{
				ID:            j,
				OP:            i,
				NaturalFreqHz: float32(m.NaturalFreqHz),
				DampedFreqHz:  float32(m.DampedFreqHz),
				DampingRatio:  float32(m.DampingRatio),
			})
		}
		results.OPs = append(results.OPs,
			OperatingPoint{
				ID:        i,
				Files:     lr.FilePaths,
				RotSpeed:  float32(lr.MBC.RotSpeed),
				WindSpeed: float32(lr.MBC.WindSpeed),
				Modes:     modes,
			},
		)
	}

	return results, nil
}

func (r *Results) Save(caseDir string) error {

	// Write results data to file
	bs, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(caseDir, "results.json"), bs, 0777)
	if err != nil {
		return err
	}

	// Loop through linearization OPs
	for _, linOP := range r.LinOPs {

		linOP.RootPath = filepath.Join(caseDir, filepath.Base(linOP.RootPath))

		// Write MBC data to file
		bs, err := json.MarshalIndent(linOP.MBC, "", "\t")
		if err != nil {
			return err
		}
		err = os.WriteFile(linOP.RootPath+"_mbc.json", bs, 0777)
		if err != nil {
			return err
		}

		// Write Eigen analysis mode results data to file
		w := &bytes.Buffer{}
		linOP.Modes.ToCSV(w)
		err = os.WriteFile(linOP.RootPath+"_modes.csv", w.Bytes(), 0777)
		if err != nil {
			return err
		}
	}

	return nil
}

func LoadResults(linDir string) (*Results, error) {

	// Create results structure
	r := Results{}

	// Load results
	bs, err := os.ReadFile(filepath.Join(linDir, "results.json"))
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bs, &r); err != nil {
		return nil, err
	}

	// Loop through linearization OPs
	for i := range r.LinOPs {

		linOP := &r.LinOPs[i]

		linOP.RootPath = filepath.Join(linDir, filepath.Base(linOP.RootPath))

		// Load MBC data
		bs, err := os.ReadFile(linOP.RootPath + "_mbc.json")
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bs, &linOP.MBC)
		if err != nil {
			return nil, err
		}

		// Write Eigen analysis mode results data to file
		modesFile := linOP.RootPath + "_modes.csv"
		bs, err = os.ReadFile(modesFile)
		if err != nil {
			return nil, err
		}
		linOP.Modes, err = lin.ReadModesCSV(bytes.NewReader(bs))
		if err != nil {
			return nil, fmt.Errorf("error parsing '%s': %w", modesFile, err)
		}
	}

	return &r, nil
}
