package lin

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"math/cmplx"
	"strconv"
	"strings"
)

type Modes []Mode

type Mode struct {
	ID             int     `json:"ID"`
	OP             int     `json:"OP"`
	MaxMod         string  `json:"MaxMod"`
	NaturalFreqRaw float64 `json:"NaturalFreqRaw"`
	NaturalFreqHz  float64 `json:"NaturalFreqHz"`
	DampedFreqRaw  float64 `json:"DampedFreqRaw"`
	DampedFreqHz   float64 `json:"DampedFreqHz"`
	DampingRatio   float64 `json:"DampingRatio"`

	EigenValue   complex128   `json:"-"`
	EigenVector  []complex128 `json:"-"`
	EigenIndices []int        `json:"-"`
}

func (ms Modes) ToCSV(w io.Writer) {
	cw := csv.NewWriter(w)
	cw.Write([]string{"ID", "OP", "MaxMod", "EVSize", "IndSize", "NaturalFreqRaw", "NaturalFreqHz",
		"DampedFreqRaw", "DampedFreqHz", "DampingRatio", "EigenValue"})
	for _, m := range ms {
		rec := []string{
			strconv.Itoa(m.ID),
			strconv.Itoa(m.OP),
			m.MaxMod,
			strconv.Itoa(len(m.EigenVector)),
			strconv.Itoa(len(m.EigenIndices)),
			strconv.FormatFloat(m.NaturalFreqRaw, 'g', -1, 64),
			strconv.FormatFloat(m.NaturalFreqHz, 'g', -1, 64),
			strconv.FormatFloat(m.DampedFreqRaw, 'g', -1, 64),
			strconv.FormatFloat(m.DampedFreqHz, 'g', -1, 64),
			strconv.FormatFloat(m.DampingRatio, 'g', -1, 64),
			strconv.FormatComplex(m.EigenValue, 'g', -1, 64),
		}
		for _, v := range m.EigenVector {
			rec = append(rec, strconv.FormatComplex(v, 'g', -1, 64))
		}
		for _, v := range m.EigenIndices {
			rec = append(rec, strconv.Itoa(v))
		}
		cw.Write(rec)
	}
	cw.Flush()
}

func ReadModesCSV(r io.Reader) (Modes, error) {

	// Read all records
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	recs, err := cr.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	// If no records return
	if len(recs) <= 1 {
		return Modes{}, nil
	}

	// Remove header record
	recs = recs[1:]

	// Loop through records and create modes
	modes := Modes{}
	for _, rec := range recs {

		evSize := 0
		fmt.Sscan(rec[3], &evSize)
		indSize := 0
		fmt.Sscan(rec[4], &indSize)

		if len(rec) < 11+evSize+indSize {
			return nil, fmt.Errorf("insufficient columns in rec")
		}

		m := Mode{
			EigenVector:  make([]complex128, evSize),
			EigenIndices: make([]int, indSize),
		}
		fmt.Sscan(rec[0], &m.ID)
		fmt.Sscan(rec[1], &m.OP)
		m.MaxMod = rec[2]
		fmt.Sscan(rec[5], &m.NaturalFreqRaw)
		fmt.Sscan(rec[6], &m.NaturalFreqHz)
		fmt.Sscan(rec[7], &m.DampedFreqRaw)
		fmt.Sscan(rec[8], &m.DampedFreqHz)
		fmt.Sscan(rec[9], &m.DampingRatio)
		fmt.Sscan(rec[10], &m.EigenValue)
		rec = rec[11:]
		for i, v := range rec[:evSize] {
			fmt.Sscan(v, &m.EigenVector[i])
		}
		rec = rec[evSize:]
		for i, v := range rec[:indSize] {
			fmt.Sscan(v, &m.EigenIndices[i])
		}

		modes = append(modes, m)
	}

	return modes, nil
}

// Filter returns true if mode should be included based on arguments.
func (m Mode) Filter(freqRangeHz [2]float64, structMax bool) bool {

	// If mode natural frequency is outside range, continue
	if m.NaturalFreqHz < freqRangeHz[0] || m.NaturalFreqHz > freqRangeHz[1] {
		return false
	}

	// If structural max is enabled, return false
	if structMax {

		// If maximum eigenvector magnitude occurs for a state not in ElastoDyn, BeamDyn, or SubDyn, return false
		if !(strings.HasPrefix(m.MaxMod, "ED") || strings.HasPrefix(m.MaxMod, "BD") || strings.HasPrefix(m.MaxMod, "SD")) {
			return false
		}
	}

	return true
}

// MAC returns the modal assurance criteria indicating mode shape similarity.
// 0=no correlation, 1=total correlation.
func (md1 Mode) MAC(md2 *Mode) (float64, error) {

	if len(md1.EigenIndices) != len(md2.EigenIndices) {
		return 0, fmt.Errorf("EigenVectors are different lengths")
	}

	var numer complex128
	var denom1, denom2 complex128
	for _, i := range md1.EigenIndices {
		numer += md1.EigenVector[i] * cmplx.Conj(md2.EigenVector[i])
		denom1 += md1.EigenVector[i] * cmplx.Conj(md1.EigenVector[i])
		denom2 += md2.EigenVector[i] * cmplx.Conj(md2.EigenVector[i])
	}

	mac := math.Pow(cmplx.Abs(numer), 2) / cmplx.Abs(denom1*denom2)

	return mac, nil
}

// https://past.isma-isaac.be/downloads/isma2010/papers/isma2010_0103.pdf
func (md1 Mode) MACX(md2 *Mode) (float64, error) {

	if len(md1.EigenIndices) != len(md2.EigenIndices) {
		return 0, fmt.Errorf("EigenVectors are different lengths")
	}

	var numer1, numer2, denom11, denom12, denom21, denom22 complex128
	for _, i := range md1.EigenIndices {
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

	if len(md1.EigenIndices) != len(md2.EigenIndices) {
		return 0, fmt.Errorf("EigenVectors are different lengths")
	}

	mu1 := md1.EigenVector
	mu2 := md2.EigenVector
	lam1 := md1.EigenValue
	lam2 := md2.EigenValue

	var numer1, numer2, denom11, denom12, denom21, denom22 complex128

	for _, i := range md1.EigenIndices {
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

func (m *Mode) Name() string {
	return fmt.Sprintf("%d-%d", m.OP, m.ID)
}
