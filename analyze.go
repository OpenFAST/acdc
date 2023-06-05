package main

type Analyze struct {
	IncludeAero   bool             `json:"IncludeAero"`
	NumCPUs       int              `json:"NumCPUs"`
	StructuralOPs []OperatingPoint `json:"StructuralOPs"`
	AeroStructOPs []OperatingPoint `json:"AeroStructOPs"`
}

func NewAnalyze() *Analyze {
	return &Analyze{NumCPUs: 1}
}

type OperatingPoint struct {
	WindSpeed  float32 `json:"WindSpeed"`  // Wind speed (m/s)
	RotorSpeed float64 `json:"RotorSpeed"` // Rotor speed in (rpm)
	BladePitch float64 `json:"BladePitch"` // Blade pitch (deg)
}
