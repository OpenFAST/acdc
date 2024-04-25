package main

import (
	"acdc/lin"
	"fmt"
	"math"
	"path/filepath"
	"regexp"
)

type Results struct {
	LinDir  string           `json:"LinDir"`
	HasWind bool             `json:"HasWind"`
	MinFreq float32          `json:"MinFreq"`
	MaxFreq float32          `json:"MaxFreq"`
	OPs     []OperatingPoint `json:"OPs"`
	LinOPs  []lin.OPResult   `json:"LinOPs"`
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

	// Initialize max rotor speed
	maxRotSpeed := 0.0

	// Extract data from linearization results
	for i, lr := range linResults {
		results.HasWind = lr.MBC.WindSpeed > 0 || results.HasWind
		modes := []Mode{}
		for j, m := range lr.EigRes.Modes {
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

		// If rotor speed is above maximum, save it
		if lr.MBC.RotSpeed > maxRotSpeed {
			maxRotSpeed = lr.MBC.RotSpeed
		}
	}

	// Calculate the recommended max diagram frequency
	results.MaxFreq = float32(math.Trunc(100*maxRotSpeed/60*15) / 100)

	return results, nil
}
