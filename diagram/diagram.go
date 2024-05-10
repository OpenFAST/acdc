package diagram

import (
	"acdc/lin"
	"encoding/json"
	"os"
	"strconv"
)

// Diagram contains data for drawing the Campbell Diagram
type Diagram struct {
	HasWind    bool      `json:"HasWind"`
	RotSpeeds  []float32 `json:"RotSpeeds"`
	WindSpeeds []float32 `json:"WindSpeeds"`
	Lines      []Line    `json:"Lines"`
}

type Options struct {
	MinFreq      float64 `json:"MinFreq"`
	MaxFreq      float64 `json:"MaxFreq"`
	Cluster      bool    `json:"Cluster"`
	FilterStruct bool    `json:"FilterStruct"`
}

type Line struct {
	ID     int     `json:"ID"`
	Label  string  `json:"Label"`
	Color  string  `json:"Color"`
	Dash   []int   `json:"Dash"`
	Hidden bool    `json:"Hidden"`
	Points []Point `json:"Points"`
}

type Point struct {
	Line          int     `json:"Line"`
	OP            int     `json:"OpPtID"`
	Mode          int     `json:"ModeID"`
	RotSpeed      float32 `json:"RotSpeed"`
	WindSpeed     float32 `json:"WindSpeed"`
	NaturalFreqHz float32 `json:"NaturalFreqHz"`
	DampedFreqHz  float32 `json:"DampedFreqHz"`
	DampingRatio  float32 `json:"DampingRatio"`
}

type ModeSet struct {
	ID        int         `json:"ID"`
	Label     string      `json:"Label"`
	Frequency [2]float64  `json:"Frequency"`
	Modes     []*lin.Mode `json:"-"`
}

type ModeIndex struct {
	OP     int     `json:"OP"`
	Mode   int     `json:"Mode"`
	Weight float64 `json:"Weight"`
}

// CampbellDiagram returns a Campbell Diagram structure from the results
func New(OPs []lin.LinOP, opts Options) (*Diagram, error) {

	// Collect operating point data
	rotSpeeds := make([]float32, len(OPs))
	windSpeeds := make([]float32, len(OPs))
	hasWind := false
	for i, linOP := range OPs {
		rotSpeeds[i] = float32(linOP.MBC.RotSpeed)
		windSpeeds[i] = float32(linOP.MBC.WindSpeed)
		hasWind = linOP.MBC.WindSpeed > 0 || hasWind
	}

	// Build mode sets based on modal assurance criteria
	modeSets, err := connectModesMAC(OPs, [2]float64{opts.MinFreq, opts.MaxFreq}, opts.FilterStruct)
	if err != nil {
		return nil, err
	}

	// Refine mode sets using spectral clustering
	if opts.Cluster {
		if err := clusterModes(OPs, modeSets); err != nil {
			return nil, err
		}
	}

	// Create diagram lines
	lines := []Line{}
	for i, ms := range modeSets {
		line := Line{
			ID:     i,
			Label:  "Line " + strconv.Itoa(i+1),
			Points: make([]Point, len(ms.Modes)),
		}
		for j, m := range ms.Modes {
			line.Points[j] = Point{
				Line:          line.ID,
				OP:            m.OP,
				Mode:          m.ID,
				RotSpeed:      float32(rotSpeeds[m.OP]),
				WindSpeed:     float32(windSpeeds[m.OP]),
				NaturalFreqHz: float32(m.NaturalFreqHz),
				DampedFreqHz:  float32(m.DampedFreqHz),
				DampingRatio:  float32(m.DampingRatio),
			}
		}
		lines = append(lines, line)
	}

	// Return the diagram
	return &Diagram{
		HasWind:    hasWind,
		RotSpeeds:  rotSpeeds,
		WindSpeeds: windSpeeds,
		Lines:      lines,
	}, nil
}

func Load(path string) (*Diagram, error) {

	d := Diagram{}
	bs, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bs, &d); err != nil {
		return nil, err
	}

	return &d, nil
}
