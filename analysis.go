package main

import (
	"fmt"
	"math"
	"sort"

	"gonum.org/v1/gonum/interp"
)

type Analysis struct {
	Cases []Case `json:"Cases"`
}

func NewAnalysis() *Analysis {
	return &Analysis{
		Cases: []Case{
			NewCase(),
		},
	}
}

type Case struct {
	ID              int         `json:"ID"`
	Name            string      `json:"Name"`
	IncludeAero     bool        `json:"IncludeAero"`
	RotorSpeedRange Range       `json:"RotorSpeedRange"`
	WindSpeedRange  Range       `json:"WindSpeedRange"`
	CutIn           float32     `json:"CutIn"`
	Rated           float32     `json:"Rated"`
	CutOut          float32     `json:"CutOut"`
	Curve           []Condition `json:"Curve"`
	OperatingPoints []Condition `json:"OperatingPoints"`
}

func NewCase() Case {
	c := Case{
		ID:              1,
		Name:            "Base",
		IncludeAero:     false,
		RotorSpeedRange: Range{Min: 1, Max: 10, Num: 5},
		WindSpeedRange:  Range{Min: 1, Max: 10, Num: 5},
		Curve: []Condition{
			{WindSpeed: 1, RotorSpeed: 1, BladePitch: 0},
			{WindSpeed: 20, RotorSpeed: 10, BladePitch: 90},
		},
		OperatingPoints: []Condition{},
	}
	c.Calculate()
	return c
}

func (c *Case) Calculate() error {

	if c.IncludeAero {
		sort.SliceStable(c.Curve, func(i, j int) bool {
			return c.Curve[i].WindSpeed < c.Curve[j].WindSpeed
		})
	} else {
		sort.SliceStable(c.Curve, func(i, j int) bool {
			return c.Curve[i].RotorSpeed < c.Curve[j].RotorSpeed
		})
	}

	// Get rotor speed and blade pitch arrays
	windSpeeds := []float64{}
	rotorSpeeds := []float64{}
	bladePitches := []float64{}
	xMap := map[float64]struct{}{}
	for i, point := range c.Curve {
		c.Curve[i].ID = i + 1
		if c.IncludeAero {
			if _, ok := xMap[point.WindSpeed]; ok {
				continue
			}
			xMap[point.WindSpeed] = struct{}{}
			windSpeeds = append(windSpeeds, point.WindSpeed)
			rotorSpeeds = append(rotorSpeeds, point.RotorSpeed)
			bladePitches = append(bladePitches, point.BladePitch)
		} else {
			if _, ok := xMap[point.RotorSpeed]; ok {
				continue
			}
			xMap[point.RotorSpeed] = struct{}{}
			rotorSpeeds = append(rotorSpeeds, point.RotorSpeed)
			bladePitches = append(bladePitches, point.BladePitch)
		}
	}

	// If no valid points in curve, reset operating points and return
	if len(bladePitches) < 2 {
		c.OperatingPoints = []Condition{}
		return nil
	}

	if c.IncludeAero {

		// Allocate operating points
		c.OperatingPoints = make([]Condition, c.WindSpeedRange.Num)

		// Create spline to interpolate rotor speed from wind speed
		var rsSpline interp.NaturalCubic
		if err := rsSpline.Fit(windSpeeds, rotorSpeeds); err != nil {
			return fmt.Errorf("error fitting cubic spline to Structure Rotor Speed and Blade Pitch: %w", err)
		}

		// Create spline to interpolate blade pitch from wind speeds
		var bpSpline interp.NaturalCubic
		if err := bpSpline.Fit(windSpeeds, bladePitches); err != nil {
			return fmt.Errorf("error fitting cubic spline to Structure Rotor Speed and Blade Pitch: %w", err)
		}

		// Calculate wind speed increment
		delta := c.WindSpeedRange.Delta()

		// Populate operating points
		for i := range c.OperatingPoints {
			windSpeed := c.WindSpeedRange.Min + delta*float64(i)
			op := &c.OperatingPoints[i]
			op.ID = i + 1
			op.WindSpeed = windSpeed
			op.RotorSpeed = rsSpline.Predict(windSpeed)
			op.BladePitch = bpSpline.Predict(windSpeed)
		}

	} else {

		// Allocate operating points
		c.OperatingPoints = make([]Condition, c.RotorSpeedRange.Num)

		// Create spline to interpolate blade pitch at given rotor speeds
		var bpSpline interp.NaturalCubic
		if err := bpSpline.Fit(rotorSpeeds, bladePitches); err != nil {
			return fmt.Errorf("error fitting cubic spline to Structure Rotor Speed and Blade Pitch: %w", err)
		}

		// Calculate rotor speed increment
		delta := c.RotorSpeedRange.Delta()

		// Populate operating points
		for i := range c.OperatingPoints {
			rotorSpeed := c.RotorSpeedRange.Min + delta*float64(i)
			op := &c.OperatingPoints[i]
			op.ID = i + 1
			op.RotorSpeed = rotorSpeed
			op.BladePitch = bpSpline.Predict(rotorSpeed)
		}
	}

	return nil
}

type Condition struct {
	ID         int     `json:"ID"`
	WindSpeed  float64 `json:"WindSpeed"`  // Wind speed (m/s)
	RotorSpeed float64 `json:"RotorSpeed"` // Rotor speed in (rpm)
	BladePitch float64 `json:"BladePitch"` // Blade pitch (deg)
}

type Range struct {
	Min float64 `json:"Min"`
	Max float64 `json:"Max"`
	Num int     `json:"Num"`
}

func (r Range) Delta() float64 {
	delta := (r.Max - r.Min) / float64(r.Num-1)
	if math.IsNaN(delta) || math.IsInf(delta, 0) {
		delta = 0
	}
	return delta
}

type Structure struct {
	RotorSpeedRange Range       `json:"RotorSpeedRange"`
	Curve           []Condition `json:"Curve"`
	OperatingPoints []Condition `json:"OperatingPoints"`
}

type AeroStructure struct {
	WindSpeedRange  Range       `json:"WindSpeedRange"`
	Curve           []Condition `json:"Curve"`
	OperatingPoints []Condition `json:"OperatingPoints"`
}

func (an *Analysis) Calculate() error {

	for i := range an.Cases {
		an.Cases[i].ID = i + 1
		if err := an.Cases[i].Calculate(); err != nil {
			return fmt.Errorf("error calculating Case %d '%s': %w",
				i, an.Cases[i].Name, err)
		}
	}

	return nil
}
