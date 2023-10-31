package mbc3

import (
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/mat"
)

type MatData struct {
	NumStep          int
	AzimuthRad       []float64 // radians
	AzimuthDeg       []float64 // degrees
	Omega            []float64
	OmegaDot         []float64
	WindSpeed        []float64
	A, B, C, D       []*mat.Dense
	OP_x, OP_u, OP_y OPSlice
	OpX              *mat.Dense
	OpXd             *mat.Dense
	OpU              *mat.Dense
	OpY              *mat.Dense
}

type OPOrder struct {
	Num         int     `json:"Num"`
	NumFixed    int     `json:"NumFixed"`
	NumRotating int     `json:"NumRotating"`
	NumTriplets int     `json:"NumTriplets"`
	Indices     []int   `json:"Indices"`
	Triplets    [][]int `json:"Triplets"`
}

// Blade regex for descriptions
var bldRe = regexp.MustCompile(`(?i)(blade\s+|blade root |PitchBearing|BD_|BD)(\d)`)

func NewOPOrder(ops OPSlice, numBlades int) OPOrder {

	// Get rotating state descriptions mapped to new index
	// also remove description text in parentheses
	opDescIndexMap := map[string]int{}
	for _, op := range ops {
		if op.IsRotFrame {
			opDescIndexMap[op.SanitizedDesc()] = op.Index
		}
	}

	// Initialize fixed, rotating, and other order slices
	fixedOrder := []int{}
	rotatingOrder := []int{}
	otherOrder := []int{}

	// Create slice of blade triplets
	triplets := [][]int{}

	// Loop through operating points
	for _, op := range ops {

		// If op is non-rotating, add index to fixed slice, continue
		if !op.IsRotFrame {
			fixedOrder = append(fixedOrder, op.Index)
			continue
		}

		// Get description without parentheses
		desc := op.SanitizedDesc()

		// If description is not in map (already used), continue
		if _, ok := opDescIndexMap[desc]; !ok {
			continue
		}

		// Find blade number text via regular expression, if no match, continue
		matches := bldRe.FindStringSubmatch(desc)
		if matches == nil {
			continue
		}

		// Initialize empty blade triplet
		triplet := []int{}

		// Lookup descriptions in map with all blade numbers
		for j := 1; j <= numBlades; j++ {

			// Create description from blade number to look up in map
			testDesc := strings.Replace(desc, matches[0], matches[1]+strconv.Itoa(j), 1)

			// If description found in map, add index to triplet, remove desc from map
			if index, ok := opDescIndexMap[testDesc]; ok {
				triplet = append(triplet, index)
				delete(opDescIndexMap, testDesc)
			}
		}

		// If triplet length is not the number of blades, then OPs do not belong to blade,
		// add triplet to other order, continue
		if len(triplet) != numBlades {
			otherOrder = append(otherOrder, triplet...)
			continue
		}

		// Triplet OPs are for blade, add to slice of triplets, and rotating order
		triplets = append(triplets, triplet)
		rotatingOrder = append(rotatingOrder, triplet...)
	}

	// Return operating point order
	return OPOrder{
		Num:         len(ops),
		NumFixed:    len(fixedOrder) + len(otherOrder),
		NumRotating: len(rotatingOrder),
		NumTriplets: len(triplets),
		Indices:     append(append(fixedOrder, otherOrder...), rotatingOrder...),
		Triplets:    triplets,
	}
}

func CombineOPOrders(opOrders ...OPOrder) OPOrder {
	opOrder := OPOrder{}
	for _, opo := range opOrders {
		opOrder.NumFixed += opo.NumFixed
		opOrder.NumRotating += opo.NumRotating
		opOrder.Indices = append(opOrder.Indices, opo.Indices...)
		opOrder.Triplets = append(opOrder.Triplets, opo.Triplets...)
	}
	opOrder.Num = len(opOrder.Indices)
	opOrder.NumTriplets = len(opOrder.Triplets)
	return opOrder
}

// NewMatData returns a structure initialized from the given linearization data
func NewMatData(lds []*LinData) *MatData {

	// Sort linearization file data by azimuth
	sort.Slice(lds, func(i, j int) bool {
		return lds[i].Azimuth < lds[j].Azimuth
	})

	// Get number of linearization files
	numStep := len(lds)

	// Get linearization data to use for initialization (max azimuth)
	ld := lds[numStep-1]

	// Create matrix data structure
	md := &MatData{
		NumStep:    numStep,
		AzimuthDeg: make([]float64, numStep),
		AzimuthRad: make([]float64, numStep),
		Omega:      make([]float64, numStep),
		OmegaDot:   make([]float64, numStep),
		WindSpeed:  make([]float64, numStep),
		OP_x:       ld.OP_x,
		OP_u:       ld.OP_u,
		OP_y:       ld.OP_y,
	}

	//--------------------------------------------------------------------------
	// General data
	//--------------------------------------------------------------------------

	// Set azimuth, omega, and wind speed
	for i, ld := range lds {
		md.AzimuthRad[i] = ld.Azimuth
		md.AzimuthDeg[i] = ld.Azimuth * 180 / math.Pi
		md.Omega[i] = ld.RotorSpeed
		md.OmegaDot[i] = 0.0
		md.WindSpeed[i] = ld.WindSpeed
	}

	//--------------------------------------------------------------------------
	// State storage
	//--------------------------------------------------------------------------

	// If number of states is greater than zero
	if len(md.OP_x) > 0 {

		// Initialize and populate X matrix from lin files
		md.OpX = mat.NewDense(numStep, len(md.OP_x), nil)
		for i, ld := range lds {
			md.OpX.SetRow(i, ld.OP_x.Values())
		}

		// Initialize and populate Xdot matrix from lin files
		md.OpXd = mat.NewDense(numStep, len(md.OP_x), nil)
		for i, ld := range lds {
			md.OpXd.SetRow(i, ld.OP_xdot.Values())
		}
	}

	//--------------------------------------------------------------------------
	// Input storage
	//--------------------------------------------------------------------------

	// If number of inputs is greater than zero
	if len(md.OP_u) > 0 {
		md.OpU = mat.NewDense(numStep, len(md.OP_u), nil)
		for i, ld := range lds {
			md.OpU.SetRow(i, ld.OP_u.Values())
		}
	}

	//--------------------------------------------------------------------------
	// Output storage
	//--------------------------------------------------------------------------

	// If number of outputs is greater than zero
	if len(md.OP_y) > 0 {
		md.OpY = mat.NewDense(numStep, len(md.OP_y), nil)
		for i, ld := range lds {
			md.OpY.SetRow(i, ld.OP_y.Values())
		}
	}

	//--------------------------------------------------------------------------
	// Matrix storage
	//--------------------------------------------------------------------------

	if len(md.OP_x) > 0 {
		md.A = make([]*mat.Dense, numStep)
		for i, ld := range lds {
			md.A[i] = mat.DenseCopyOf(ld.A) // (NumX,NumX)
		}
	}
	if len(md.OP_x) > 0 && len(md.OP_u) > 0 {
		md.B = make([]*mat.Dense, numStep)
		for i, ld := range lds {
			md.B[i] = mat.DenseCopyOf(ld.B) // (NumX,NumU)
		}
	}
	if len(md.OP_y) > 0 && len(md.OP_x) > 0 {
		md.C = make([]*mat.Dense, numStep)
		for i, ld := range lds {
			md.C[i] = mat.DenseCopyOf(ld.C) // (NumY,NumX)
		}
	}
	if len(md.OP_y) > 0 && len(md.OP_u) > 0 {
		md.D = make([]*mat.Dense, numStep)
		for i, ld := range lds {
			md.D[i] = mat.DenseCopyOf(ld.D) // (NumY,NumU)
		}
	}

	return md
}
