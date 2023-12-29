package diagram

import (
	"acdc/lin"
	"fmt"
	"math"
	"sort"
	"strconv"

	"github.com/mkmik/argsort"
	"github.com/parallelo-ai/kmeans"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

// Diagram contains data for drawing the Campbell Diagram
type Diagram struct {
	HasWind    bool      `json:"HasWind"`
	RotSpeeds  []float32 `json:"RotSpeeds"`
	WindSpeeds []float32 `json:"WindSpeeds"`
	Lines      []Line    `json:"Lines"`
}

type Line struct {
	ID     int     `json:"ID"`
	Label  string  `json:"Label"`
	Points []Point `json:"Points"`
	Hide   bool    `json:"Hide"`
}

type Point struct {
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
func New(OPs []lin.OPResults, maxFreqHz float64, doCluster bool) (*Diagram, error) {

	// Collect operating point data
	rotSpeeds := make([]float32, len(OPs))
	windSpeeds := make([]float32, len(OPs))
	hasWind := false
	for i, g := range OPs {
		rotSpeeds[i] = float32(g.MBC.RotSpeed)
		windSpeeds[i] = float32(g.MBC.WindSpeed)
		hasWind = g.MBC.WindSpeed > 0 || hasWind

	}

	// Build mode sets based on modal assurance criteria
	modeSets, err := connectModesMAC(OPs, maxFreqHz)
	if err != nil {
		return nil, err
	}

	// Refine mode sets using spectral clustering
	if doCluster {
		if err := clusterModes(OPs, modeSets); err != nil {
			return nil, err
		}
	}

	// Create diagram lines
	lines := make([]Line, len(modeSets))
	for i, ms := range modeSets {
		points := make([]Point, len(ms.Modes))
		for j, m := range ms.Modes {
			points[j] = Point{
				OP:            m.OP,
				Mode:          m.ID,
				RotSpeed:      float32(rotSpeeds[m.OP]),
				WindSpeed:     float32(windSpeeds[m.OP]),
				NaturalFreqHz: float32(m.NaturalFreqHz),
				DampedFreqHz:  float32(m.DampedFreqHz),
				DampingRatio:  float32(m.DampingRatio),
			}
		}
		lines[i] = Line{
			ID:     i + 1,
			Label:  strconv.Itoa(i + 1),
			Points: points,
		}
	}

	// Return the diagram
	return &Diagram{
		HasWind:    hasWind,
		RotSpeeds:  rotSpeeds,
		WindSpeeds: windSpeeds,
		Lines:      lines,
	}, nil
}

// connectModesMAC builds connected sets of modes from linearization results
func connectModesMAC(OPs []lin.OPResults, maxFreqHz float64) ([]*ModeSet, error) {

	// Create array of mode sets
	modeSets := []*ModeSet{}

	// Loop through modes in first operating point and add to map
	for i := range OPs[0].EigRes.Modes {

		// Get pointer to mode
		m := &OPs[0].EigRes.Modes[i]

		// If mode natural frequency exceeds limit, continue
		if m.NaturalFreqHz > maxFreqHz {
			continue
		}

		// Initialize mode set with mode
		modeSets = append(modeSets, &ModeSet{
			ID:    i,
			Label: fmt.Sprintf("%d", i),
			Modes: []*lin.Mode{m},
		})
	}

	// Loop through operating point results
	for opID, op := range OPs {

		// Skip first operating point
		if opID == 0 {
			continue
		}

		// Create empty weighting matrix
		w := mat.NewDense(len(modeSets), len(op.EigRes.Modes), nil)

		// Create map mapping mode index to mode
		modeIndexMap := map[int]*lin.Mode{}

		// Loop through modes in mode set map
		for j, modeSet := range modeSets {

			// 	Get last mode in mode set
			mp := modeSet.Modes[len(modeSet.Modes)-1]

			// Loop through modes in current operating point
			k := 0
			for l := range op.EigRes.Modes {

				// Get mode
				mn := &op.EigRes.Modes[l]

				// If mode natural frequency exceeds limit, continue
				if mn.NaturalFreqHz > maxFreqHz {
					continue
				}

				// Calculate MAC between modes
				mac, err := mp.MAC(mn)
				if err != nil {
					return nil, err
				}

				// Modify MAC by change in frequency
				mac *= 1 - math.Abs(mn.NaturalFreqHz-mp.NaturalFreqHz)/maxFreqHz

				// Add MAC to weight matrix
				w.Set(j, k, mac)

				// Add mode to index map
				modeIndexMap[k] = mn

				k++
			}
		}

		// Get max weight value
		wMax := mat.Max(w)

		// Create cost matrix (ints) from weights (rescale to maximize precision)
		cost := NewIntMatrix(len(modeSets), len(modeIndexMap), 0)
		for j := range cost {
			for k := range cost[j] {
				v := w.At(j, k)
				cost[j][k] = int(1e7 * (1 - v/wMax))
			}
		}

		// Find mode pairings that minimizes the total cost
		pairs, err := MinCostAssignment(cost)
		if err != nil {
			return nil, err
		}

		// Add connected modes to sets
		for _, pair := range pairs {

			// Look up mode set from previous mode index
			modeSet := modeSets[pair[0]]

			// Add paired mode to slice of modes
			modeSet.Modes = append(modeSet.Modes, modeIndexMap[pair[1]])

			// Remove paired mode from map
			delete(modeIndexMap, pair[1])
		}

		// Loop through unpaired modes and create new mode sets
		for _, m := range modeIndexMap {
			modeSets = append(modeSets, &ModeSet{
				ID:    len(modeSets),
				Label: fmt.Sprintf("%d", len(modeSets)),
				Modes: []*lin.Mode{m},
			})
		}
	}

	// Create temporary slice for filtering mode sets
	allModeSets := modeSets
	modeSets = modeSets[:0]

	// Loop through mode sets
	for _, modeSet := range allModeSets {

		// Skip empty mode sets
		if len(modeSet.Modes) == 0 {
			continue
		}

		// Get min and max frequency from first mode in set
		m := modeSet.Modes[0]
		modeSet.Frequency = [2]float64{m.NaturalFreqHz, m.NaturalFreqHz}

		// Calculate min and max natural frequency from remaining indices
		for _, m := range modeSet.Modes[1:] {
			modeSet.Frequency[0] = min(modeSet.Frequency[0], m.NaturalFreqHz)
			modeSet.Frequency[1] = max(modeSet.Frequency[1], m.NaturalFreqHz)
		}

		// Append mode set to results
		modeSets = append(modeSets, modeSet)
	}

	// Sort mode sets by minimum frequency
	sort.Slice(modeSets, func(i, j int) bool {
		return modeSets[i].Frequency[0] < modeSets[j].Frequency[0]
	})

	return modeSets, nil
}

func clusterModes(OPs []lin.OPResults, modeSets []*ModeSet) error {

	// Find groups of potentially overlapping mode sets
	modeSetGroups := [][]*ModeSet{{}}
	j := 0
	for _, ms := range modeSets {

		// If mode set is incomplete, fewer indices than OPs, don't include in group
		if len(ms.Modes) < len(OPs) {
			continue
		}

		opMap1 := map[int]*lin.Mode{}
		for _, m := range ms.Modes {
			opMap1[m.OP] = m
		}

		// Find minimum gap compared to any mode set in group
		minGap := 1000.0
		for _, msg := range modeSetGroups[j] {
			opMap2 := map[int]*lin.Mode{}
			for _, m := range msg.Modes {
				opMap2[m.OP] = m
			}
			for i := 0; i < len(OPs); i++ {
				m1, ok1 := opMap1[i]
				m2, ok2 := opMap2[i]
				if ok1 && ok2 {
					minGap = min(m1.NaturalFreqHz-m2.NaturalFreqHz, minGap)
				}
			}
		}

		if len(modeSetGroups[j]) > 0 && minGap > 0.05 {
			modeSetGroups = append(modeSetGroups, []*ModeSet{})
			j++
		}

		// Add set to last group
		modeSetGroups[j] = append(modeSetGroups[j], ms)
	}

	// Loop through mode set groups. If more than one mode set in group,
	// perform spectral clustering to identify shared modes
	for _, group := range modeSetGroups {
		if len(group) > 1 {
			if err := spectralClustering(group); err != nil {
				return err
			}
		}
	}

	return nil
}

func spectralClustering(modeSets []*ModeSet) error {

	// Collect all modes in mode sets
	modes := []*lin.Mode{}
	for _, ms := range modeSets {
		modes = append(modes, ms.Modes...)
	}

	N := len(modes)

	// Create weight matrix by comparing modes with MAC
	W := mat.NewDense(N, N, nil)
	D := mat.NewDense(N, N, nil)
	D_isr := mat.NewDense(N, N, nil)
	for i, m1 := range modes {
		for j, m2 := range modes {
			if i != j {
				if mac, err := m1.MAC(m2); err == nil {
					W.Set(i, j, mac)
				}
			}
		}
		di := mat.Sum(W.RowView(i))
		D.Set(i, i, di)
		D_isr.Set(i, i, 1/math.Sqrt(di))
	}

	// Calculate Laplacian matrix (D - W)
	L := mat.NewDense(N, N, nil)
	L.Sub(D, W)

	// Calculate the symmetric laplacian (Lsym = D^{-1/2}*L*D^{-1/2})
	Lsym := mat.NewDense(N, N, nil)
	Lsym.Mul(D_isr, L)
	Lsym.Mul(Lsym, D_isr)

	// Calculate eigenvalues and eigenvectors of Lsym matrix
	eig := mat.Eigen{}
	if ok := eig.Factorize(Lsym, mat.EigenRight); !ok {
		return fmt.Errorf("error computing eigenvalues")
	}
	eigenValues := eig.Values(nil)
	eigenVectors := &mat.CDense{}
	eig.VectorsTo(eigenVectors)

	// Get indices that would sort from largest to smallest eigenvalues
	indices := argsort.SortSlice(eigenValues, func(i, j int) bool {
		return real(eigenValues[i]) < real(eigenValues[j])
	})

	numClusters := len(modeSets)
	numDims := min(int(math.Ceil(1*float64(numClusters))), N)

	d := make(kmeans.Observations, N)
	for i := 0; i < N; i++ {
		row := make([]float64, numDims)
		for j, ind := range indices[:numDims] {
			row[j] = real(eigenVectors.At(i, ind))
		}
		floats.Scale(1/floats.Norm(row, 2), row)
		d[i] = Observation(row)
	}

	// Create KMeans object with options
	km, err := kmeans.NewKmeansWithOptions(0.001, 1000)
	if err != nil {
		return err
	}

	clusterModesMap := map[int][]*lin.Mode{}
	modeClusterMap := map[*lin.Mode]int{}
	minRepeatedModes := N

	for i := 0; i < 1000; i++ {

		// Partition the data points
		clusters, err := km.Partition(d, numClusters, 0)
		if err != nil {
			return err
		}

		// Get cluster number for each mode (starts at 0)
		localClusterModesMap := map[int][]*lin.Mode{}
		localModeClusterMap := map[*lin.Mode]int{}
		for i, obs := range d {
			c := clusters.Nearest(obs)
			if _, ok := localClusterModesMap[c]; !ok {
				localClusterModesMap[c] = []*lin.Mode{}
			}
			localClusterModesMap[c] = append(localClusterModesMap[c], modes[i])
			localModeClusterMap[modes[i]] = c
		}

		// Calculate number of modes in same OP across clusters.
		// This represents how well the kmeans was able to redraw the paths
		// in the mode sets
		numRepeatedModes := 0
		for _, modes := range localClusterModesMap {
			opMap := map[int]int{}
			for _, m := range modes {
				opMap[m.OP]++
			}
			for _, numModes := range opMap {
				numRepeatedModes += numModes - 1
			}
		}

		if numRepeatedModes < minRepeatedModes {
			minRepeatedModes = numRepeatedModes
			clusterModesMap = localClusterModesMap
			modeClusterMap = localModeClusterMap
		}

		// If number of repeated OPs is sufficiently small, compared to the
		// number of modes, exit loop
		if ratio := float64(numRepeatedModes) / float64(N); ratio < 0.01 {
			break
		}
	}

	// Build cost matrix for determining which mode set goes with each cluster
	C := make([][]int, len(modeSets))
	for i := range C {
		C[i] = make([]int, numClusters)
		for _, m := range modeSets[i].Modes {
			C[i][modeClusterMap[m]]++
		}
		for j := range C[i] {
			C[i][j] = N - 1 - C[i][j]
		}
	}

	// Pair mode set with cluster which have the most overlapping modes
	modeSetClusterPairs, err := MinCostAssignment(C)
	if err != nil {
		return err
	}

	// Loop through mode set -> cluster pairings
	for _, pair := range modeSetClusterPairs {

		// Get mode set index and cluster number
		msi, cn := pair[0], pair[1]

		// Get mode set and modes in cluster
		ms := modeSets[msi]
		modes := clusterModesMap[cn]

		// Collect cluster modes by operating point
		opModesMap := map[int][]*lin.Mode{}
		for _, m := range modes {
			opModesMap[m.OP] = append(opModesMap[m.OP], m)
		}

		// Rebuild list of modes keeping one for each OP
		modes = []*lin.Mode{}
		for _, opModes := range opModesMap {
			if len(opModes) == 1 {
				modes = append(modes, opModes...)
			}
		}
		sort.Slice(modes, func(i, j int) bool {
			return modes[i].OP < modes[j].OP
		})

		// Reset modes slice and add modes which only have one OP
		ms.Modes = modes
	}

	return nil
}

type Observation []float64

func (obs Observation) Coordinates() kmeans.Coordinates {
	return kmeans.Coordinates(obs)
}

func (obs Observation) Distance(point kmeans.Coordinates) float64 {
	diff := make([]float64, len(obs))
	floats.SubTo(diff, obs, []float64(point))
	return floats.Norm(diff, 2)
}
