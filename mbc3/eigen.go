package mbc3

import (
	"fmt"
	"math"
	"math/cmplx"

	"gonum.org/v1/gonum/mat"
)

type EigenResults struct {
	Modes []ModeResults
}

type ModeResults struct {
	EigenValueReal float64
	EigenValueImag float64
	NaturalFreqRaw float64
	NaturalFreqHz  float64
	DampedFreqRaw  float64
	DampedFreqHz   float64
	DampingRatio   float64
	Magnitudes     []float64
	Phases         []float64
	Shape          []float64
}

func EigenAnalysis(A *mat.Dense, rows []int) (*EigenResults, error) {

	// Calculate eigenvalues/eigenvectors analysis
	eig := mat.Eigen{}
	if ok := eig.Factorize(A, mat.EigenRight); !ok {
		return nil, fmt.Errorf("error computing eigenvalues")
	}
	eigenVectors := &mat.CDense{}
	eig.VectorsTo(eigenVectors)

	// Create slice of mode results
	modes := []ModeResults{}

	// Collect mode results
	for i, ev := range eig.Values(nil) {

		// Skip negative imaginary eigenvalues
		if imag(ev) <= 0 {
			continue
		}

		// Create mode
		mode := ModeResults{
			EigenValueReal: real(ev),
			EigenValueImag: imag(ev),
			NaturalFreqRaw: cmplx.Abs(ev),
			NaturalFreqHz:  cmplx.Abs(ev) / (2 * math.Pi),
			DampedFreqRaw:  imag(ev),
			DampedFreqHz:   imag(ev) / (2 * math.Pi),
			DampingRatio:   -real(ev) / cmplx.Abs(ev),
			Magnitudes:     make([]float64, len(rows)),
			Phases:         make([]float64, len(rows)),
			Shape:          make([]float64, len(rows)),
		}

		// Convert eigenvector value to magnitude and phase
		for j, r := range rows {
			mode.Magnitudes[j], mode.Phases[j] = cmplx.Polar(eigenVectors.At(r, i))
		}

		// Normalize magnitudes to get mode shape
		maxMag := mode.Magnitudes[0]
		for _, m := range mode.Magnitudes {
			if math.Abs(m) > math.Abs(maxMag) {
				maxMag = m
			}
		}
		for j, m := range mode.Magnitudes {
			mode.Shape[j] = m / maxMag
		}

		// Add mode to slice of modes
		modes = append(modes, mode)
	}

	return &EigenResults{Modes: modes}, nil
}
