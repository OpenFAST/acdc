package mbc3

import (
	"fmt"
	"math"
	"math/cmplx"
	"sort"

	"gonum.org/v1/gonum/mat"
)

type Mode struct {
	ID             int          `json:"ID"`
	OP             int          `json:"OP"`
	EigenValueReal float64      `json:"EigenValueReal"`
	EigenValueImag float64      `json:"EigenValueImag"`
	NaturalFreqRaw float64      `json:"NaturalFreqRaw"`
	NaturalFreqHz  float64      `json:"NaturalFreqHz"`
	DampedFreqRaw  float64      `json:"DampedFreqRaw"`
	DampedFreqHz   float64      `json:"DampedFreqHz"`
	DampingRatio   float64      `json:"DampingRatio"`
	Magnitudes     []float64    `json:"Magnitudes"`
	Phases         []float64    `json:"Phases"`
	EigenValue     complex128   `json:"-"`
	EigenVector    []complex128 `json:"-"`
	Cluster        int          `json:"Cluster"`
}

// MAC returns the modal assurance criteria indicating mode shape similarity.
// 0=no correlation, 1=total correlation.
func (md1 Mode) MAC(md2 Mode) (float64, error) {

	if len(md1.EigenVector) != len(md2.EigenVector) {
		return 0, fmt.Errorf("EigenVectors are different lengths")
	}

	var numer complex128
	var denom1, denom2 complex128
	for i := range md1.EigenVector {
		numer += md1.EigenVector[i] * cmplx.Conj(md2.EigenVector[i])
		denom1 += md1.EigenVector[i] * cmplx.Conj(md1.EigenVector[i])
		denom2 += md2.EigenVector[i] * cmplx.Conj(md2.EigenVector[i])
	}

	mac := math.Pow(cmplx.Abs(numer), 2) / cmplx.Abs(denom1*denom2)

	return mac, nil
}

// https://past.isma-isaac.be/downloads/isma2010/papers/isma2010_0103.pdf
func (md1 Mode) MACX(md2 Mode) (float64, error) {

	if len(md1.EigenVector) != len(md2.EigenVector) {
		return 0, fmt.Errorf("EigenVectors are different lengths")
	}

	var numer1, numer2, denom11, denom12, denom21, denom22 complex128
	for i := range md1.EigenVector {
		numer1 += md1.EigenVector[i] * cmplx.Conj(md2.EigenVector[i])
		numer2 += md1.EigenVector[i] * md2.EigenVector[i]

		denom11 += md1.EigenVector[i] * cmplx.Conj(md1.EigenVector[i])
		denom12 += md2.EigenVector[i] * md2.EigenVector[i]

		denom21 += md2.EigenVector[i] * cmplx.Conj(md2.EigenVector[i])
		denom22 += md2.EigenVector[i] * md2.EigenVector[i]
	}

	mac := math.Pow(cmplx.Abs(numer1)+cmplx.Abs(numer2), 2) /
		((cmplx.Abs(denom11) + cmplx.Abs(denom12)) *
			(cmplx.Abs(denom21) + cmplx.Abs(denom22)))

	return mac, nil
}

// https://past.isma-isaac.be/downloads/isma2010/papers/isma2010_0103.pdf
func (md1 *Mode) MACXP(md2 *Mode) (float64, error) {

	if len(md1.EigenVector) != len(md2.EigenVector) {
		return 0, fmt.Errorf("EigenVectors are different lengths")
	}

	mu1 := md1.EigenVector
	mu2 := md2.EigenVector
	lam1 := md1.EigenValue
	lam2 := md2.EigenValue

	var numer1, numer2, denom11, denom12, denom21, denom22 complex128
	for i := range mu1 {
		numer1 += cmplx.Conj(mu1[i]) * mu2[i]
		numer2 += mu1[i] * mu2[i]

		denom11 += cmplx.Conj(mu1[i]) * mu1[i]
		denom12 += mu1[i] * mu1[i]

		denom21 += cmplx.Conj(mu2[i]) * mu2[i]
		denom22 += mu2[i] * mu2[i]
	}

	num := cmplx.Abs(numer1)/cmplx.Abs(cmplx.Conj(lam1)+lam2) +
		cmplx.Abs(numer2)/cmplx.Abs(lam1+lam2)

	den := (cmplx.Abs(denom11)/(2*real(lam1)) + cmplx.Abs(denom12)/(2*cmplx.Abs(lam1))) *
		(cmplx.Abs(denom21)/(2*real(lam2)) + cmplx.Abs(denom22)/(2*cmplx.Abs(lam2)))

	return math.Pow(num, 2) / den, nil
}

func (mbc MBC) EigenAnalysis() ([]Mode, error) {

	// Get indices of eigenvector rows to keep
	rows := mbc.OrderEigen.Indices

	// Calculate eigenvalues/eigenvectors analysis
	eig := mat.Eigen{}
	if ok := eig.Factorize(mbc.AvgA, mat.EigenRight); !ok {
		return nil, fmt.Errorf("error computing eigenvalues")
	}
	eigenVectors := &mat.CDense{}
	eig.VectorsTo(eigenVectors)

	// Create slice of mode results
	modes := []Mode{}

	// Collect mode results
	for i, ev := range eig.Values(nil) {

		// Skip negative imaginary eigenvalues
		if imag(ev) <= 0 {
			continue
		}

		// Create mode
		mode := Mode{
			EigenValueReal: real(ev),
			EigenValueImag: imag(ev),
			NaturalFreqRaw: cmplx.Abs(ev),
			NaturalFreqHz:  cmplx.Abs(ev) / (2 * math.Pi),
			DampedFreqRaw:  imag(ev),
			DampedFreqHz:   imag(ev) / (2 * math.Pi),
			DampingRatio:   -real(ev) / cmplx.Abs(ev),
			Magnitudes:     make([]float64, len(rows)),
			Phases:         make([]float64, len(rows)),
			EigenValue:     ev,
			EigenVector:    make([]complex128, len(rows)),
		}

		// Store eigenvector for given rows
		for j, r := range rows {
			mode.EigenVector[j] = eigenVectors.At(r, i)
			mode.Magnitudes[j], mode.Phases[j] = cmplx.Polar(mode.EigenVector[j])
		}

		// Add mode to slice of modes
		modes = append(modes, mode)
	}

	// Sort modes by natural frequency, ascending
	sort.Slice(modes, func(i, j int) bool {
		return modes[i].NaturalFreqRaw < modes[j].NaturalFreqRaw
	})

	// Set mode identifiers
	for i := range modes {
		modes[i].ID = i + 1
	}

	return modes, nil
}
