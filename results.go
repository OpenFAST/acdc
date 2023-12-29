package main

import (
	"acdc/lin"
	"fmt"
	"path/filepath"
	"regexp"
)

type Results struct {
	LinDir  string           `json:"LinDir"`
	HasWind bool             `json:"HasWind"`
	OPs     []OperatingPoint `json:"OPs"`
	LinOPs  []lin.OPResults  `json:"LinOPs"`
}

type OperatingPoint struct {
	ID        int      `json:"ID"`
	Files     []string `json:"Files"`
	RotSpeed  float64  `json:"RotSpeed"`  // RPM
	WindSpeed float64  `json:"WindSpeed"` // m/s
	Modes     []Mode   `json:"Modes"`
}

type Mode struct {
	ID            int     `json:"ID"`
	OP            int     `json:"OP"`
	NaturalFreqHz float64 `json:"NaturalFreqHz"`
	DampedFreqHz  float64 `json:"DampedFreqHz"`
	DampingRatio  float64 `json:"DampingRatio"`
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
		for j, m := range lr.EigRes.Modes {
			modes = append(modes, Mode{
				ID:            j,
				OP:            i,
				NaturalFreqHz: m.NaturalFreqHz,
				DampedFreqHz:  m.DampedFreqHz,
				DampingRatio:  m.DampingRatio,
			})
		}
		results.OPs = append(results.OPs,
			OperatingPoint{
				ID:        i,
				RotSpeed:  lr.MBC.RotSpeed,
				WindSpeed: lr.MBC.WindSpeed,
				Modes:     modes,
			},
		)
	}

	return results, nil
}
