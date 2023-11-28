package main

import (
	"acdc/mbc3"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/mkmik/argsort"
	"github.com/parallelo-ai/kmeans"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

type Results struct {
	OPs      []OPResults        `json:"OPs"`
	ModeSets []ModeSet          `json:"ModeSets"`
	CD       CampbellDiagram    `json:"CD"`
	MBC      []*mbc3.MBC        `json:"-"`
	EigRes   *mbc3.EigenResults `json:"-"`
	Options  ResultOptions      `json:"Options"`
}

type ResultOptions struct {
	ModeScale float32 `json:"ModeScale"`
}

type OPResults struct {
	ID        int         `json:"ID"`
	Files     []string    `json:"Files"`
	RotSpeed  float64     `json:"RotSpeed"`  // RPM
	WindSpeed float64     `json:"WindSpeed"` // m/s
	Modes     []mbc3.Mode `json:"Modes"`
	Costs     [][]int     `json:"Costs"`
}

type ModeSet struct {
	ID        int          `json:"ID"`
	Label     string       `json:"Label"`
	Frequency [2]float64   `json:"Frequency"`
	Modes     []*mbc3.Mode `json:"Modes"`
}

type ModeIndex struct {
	OP     int     `json:"OP"`
	Mode   int     `json:"Mode"`
	Weight float64 `json:"Weight"`
}

func (res *Results) ProcessFiles(LinFiles []string) (err error) {

	// Organize linearization files by operating point
	linFileGroups := map[string][]string{}
	for _, filePath := range LinFiles {
		tmp := strings.Split(filePath, ".")
		filePathNoLinExt := strings.Join(tmp[:len(tmp)-2], ".")
		fileSlice := linFileGroups[filePathNoLinExt]
		fileSlice = append(fileSlice, filePath)
		linFileGroups[filePathNoLinExt] = fileSlice
	}

	linFileGroupChan := make(chan LinFileGroup, len(linFileGroups))
	linFileResultChan := make(chan LinFileResult, len(linFileGroups))

	// Launch workers
	for i := 0; i < min(len(linFileGroups), 1+2*runtime.NumCPU()/3); i++ {
		go linFileWorker(linFileGroupChan, linFileResultChan)
	}

	// Distribute linearization file groups to workers
	for groupName, files := range linFileGroups {
		linFileGroupChan <- LinFileGroup{Name: groupName, Files: files}
	}

	// Close group chan
	close(linFileGroupChan)

	// Collect results
	linFileResults := make([]LinFileResult, len(linFileGroups))
	for i := 0; i < len(linFileGroups); i++ {
		linFileResults[i] = <-linFileResultChan
	}

	for i := range linFileResults {
		if linFileResults[i].err != nil {
			return linFileResults[i].err
		}
	}

	// Sort results by group name
	sort.Slice(linFileResults, func(i, j int) bool {
		return linFileResults[i].Name < linFileResults[j].Name
	})

	// Reset members
	res.OPs = res.OPs[:0]
	res.MBC = res.MBC[:0]

	// Loop through results and add
	for i, lfr := range linFileResults {

		// Set operating point identifier for modes
		for j := range lfr.EigRes.Modes {
			lfr.EigRes.Modes[j].OP = i
		}

		// Store data in results
		res.MBC = append(res.MBC, lfr.MBC)
		res.OPs = append(res.OPs, OPResults{
			ID:        i,
			RotSpeed:  lfr.MBC.RotSpeed,
			WindSpeed: lfr.MBC.WindSpeed,
			Modes:     lfr.EigRes.Modes,
		})
	}

	// Identify modes
	err = res.BuildModeSets()
	if err != nil {
		return err
	}

	// If no mode sets found, return
	if len(res.ModeSets) == 0 {
		return nil
	}

	// Find groups of potentially overlapping mode sets
	modeSetGroups := [][]*ModeSet{{}}
	j := 0
	for i := range res.ModeSets {

		// Get pointer to the mode set
		ms := &res.ModeSets[i]

		// If mode set is incomplete, fewer indices than OPs, don't include in group
		if len(ms.Modes) < len(res.OPs) {
			continue
		}

		opMap1 := map[int]*mbc3.Mode{}
		for _, m := range ms.Modes {
			opMap1[m.OP] = m
		}

		// Find minimum gap compared to any mode set in group
		minGap := 1000.0
		for _, msg := range modeSetGroups[j] {
			opMap2 := map[int]*mbc3.Mode{}
			for _, m := range msg.Modes {
				opMap2[m.OP] = m
			}
			for i := 0; i < len(res.OPs); i++ {
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
	for _, msg := range modeSetGroups {
		if len(msg) > 1 {
			err = res.SpectralClustering(msg)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type LinFileGroup struct {
	Name  string
	Files []string
}

type LinFileResult struct {
	Name   string
	MBC    *mbc3.MBC
	EigRes *mbc3.EigenResults
	err    error
}

func linFileWorker(linFilesChan <-chan LinFileGroup, resultsChan chan<- LinFileResult) {

	// Loop through linearization files sent on channel
	for linFileGroup := range linFilesChan {

		if len(linFileGroup.Files) == 0 {
			continue
		}

		// Read linearization files
		linFileData := make([]*mbc3.LinData, len(linFileGroup.Files))
		var err error
		for i, linFilePath := range linFileGroup.Files {
			linFileData[i], err = mbc3.ReadLinFile(linFilePath)
			if err != nil {
				resultsChan <- LinFileResult{err: err}
				return
			}
		}

		// Extract matrix data from linearization file data
		matData := mbc3.NewMatData(linFileData)

		// Perform multi-blade coordinate transform
		mbc, err := matData.MBC3()
		if err != nil {
			resultsChan <- LinFileResult{err: err}
			return
		}

		// Perform Eigenanalysis to get modes
		eigRes, err := mbc.EigenAnalysis()
		if err != nil {
			resultsChan <- LinFileResult{err: err}
			return
		}

		// Write MBC data to file
		bs, err := json.MarshalIndent(mbc, "", "\t")
		if err != nil {
			resultsChan <- LinFileResult{err: err}
			return
		}
		err = os.WriteFile(linFileGroup.Name+"_mbc.json", bs, 0777)
		if err != nil {
			resultsChan <- LinFileResult{err: err}
			return
		}

		// Write Eigen analysis results data to file
		bs, err = json.MarshalIndent(eigRes.Modes, "", "\t")
		if err != nil {
			resultsChan <- LinFileResult{err: err}
			return
		}
		err = os.WriteFile(linFileGroup.Name+"_modes.json", bs, 0777)
		if err != nil {
			resultsChan <- LinFileResult{err: err}
			return
		}

		// Send MBC and mode results
		resultsChan <- LinFileResult{
			Name:   linFileGroup.Name,
			MBC:    mbc,
			EigRes: eigRes,
		}
	}
}

func (r *Results) BuildModeSets() error {

	// Set max frequency to consider
	maxRotSpeed := 0.0
	for _, op := range r.OPs {
		maxRotSpeed = max(maxRotSpeed, op.RotSpeed)
	}
	MaxFreqHz := maxRotSpeed / 60 * 18

	// Create a map of mode sets
	modeSets := []*ModeSet{}

	// Loop through modes in first operating point
	for i := range r.OPs[0].Modes {

		m := &r.OPs[0].Modes[i]

		// If mode natural frequency exceeds limit, continue
		if m.NaturalFreqHz > MaxFreqHz {
			continue
		}

		// Initialize mode set with mode
		modeSets = append(modeSets, &ModeSet{
			ID:    i,
			Label: fmt.Sprintf("%d", i),
			Modes: []*mbc3.Mode{m},
		})
	}

	// Loop through operating point results
	for opID, op := range r.OPs {

		// Skip first operating point
		if opID == 0 {
			continue
		}

		// Create weighting matrix
		w := mat.NewDense(len(modeSets), len(op.Modes), nil)

		modeIndexMap := map[int]*mbc3.Mode{}

		// Loop through modes in mode set map
		for j, modeSet := range modeSets {

			// 	Get last mode in mode set
			mp := modeSet.Modes[len(modeSet.Modes)-1]

			// Loop through modes in current operating point
			k := 0
			for l := range op.Modes {

				mn := &op.Modes[l]

				// If mode natural frequency exceeds limit, continue
				if mn.NaturalFreqHz > MaxFreqHz {
					continue
				}

				// Calculate MAC between modes
				mac, err := mp.MAC(mn)
				if err != nil {
					return err
				}

				mac *= 1 - math.Abs(mn.NaturalFreqHz-mp.NaturalFreqHz)/MaxFreqHz

				// Add MAC to weight matrix
				w.Set(j, k, mac)

				// Add mode to index map
				modeIndexMap[k] = mn

				k++
			}
		}

		// Get max weight value
		wMax := mat.Max(w)

		// Create cost matrix from weights (rescale to maximize precision)
		cost := NewIntMatrix(len(modeSets), len(modeIndexMap), 0)
		for j := range cost {
			for k := range cost[j] {
				v := w.At(j, k)
				cost[j][k] = int(1e7 * (1 - v/wMax))
			}
		}

		// Save cost matrix in operating point
		r.OPs[opID].Costs = cost

		// Find mode pairings that minimizes the total cost
		pairs, err := MinCostAssignment(cost)
		if err != nil {
			return err
		}

		// Set mode connections
		for _, pair := range pairs {

			// Look up mode set from previous mode index
			modeSet := modeSets[pair[0]]

			// Add paired mode to slice of modes
			modeSet.Modes = append(modeSet.Modes, modeIndexMap[pair[1]])
		}

		// TODO: add logic for adding new mode sets outside pairs
	}

	// Clear mode sets in results
	r.ModeSets = []ModeSet{}

	// Loop through mode sets
	for _, modeSet := range modeSets {

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
		r.ModeSets = append(r.ModeSets, *modeSet)
	}

	// Sort mode sets by minimum frequency
	sort.Slice(r.ModeSets, func(i, j int) bool {
		return r.ModeSets[i].Frequency[0] < r.ModeSets[j].Frequency[0]
	})

	return nil
}

func (r *Results) SpectralClustering(modeSets []*ModeSet) error {

	// Collect all modes in mode sets
	modes := []*mbc3.Mode{}
	for _, ms := range modeSets {
		modes = append(modes, ms.Modes...)
	}

	N := len(modes)

	// Create weight matrix by comparing modes with MAC
	W := mat.NewDense(N, N, nil)
	D := mat.NewDense(N, N, nil)
	Disr := mat.NewDense(N, N, nil)
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
		Disr.Set(i, i, 1/math.Sqrt(di))
	}

	// Calculate Laplacian matrix (D - W)
	L := mat.NewDense(N, N, nil)
	L.Sub(D, W)

	// Calculate the symmetric laplacian (Lsym = D^{-1/2}*L*D^{-1/2})
	Lsym := mat.NewDense(N, N, nil)
	Lsym.Mul(Disr, L)
	Lsym.Mul(Lsym, Disr)

	// Calculate eigenvalues and eigenvectors of Lsym matrix
	eig := mat.Eigen{}
	if ok := eig.Factorize(Lsym, mat.EigenRight); !ok {
		return fmt.Errorf("error computing eigenvalues")
	}
	eigenValues := eig.Values(nil)
	eigenVectors := &mat.CDense{}
	eig.VectorsTo(eigenVectors)

	// Get indices that would sort from largest to smallest eigenvalues
	inds := argsort.SortSlice(eigenValues, func(i, j int) bool {
		return real(eigenValues[i]) < real(eigenValues[j])
	})

	numClusters := len(modeSets)
	numDims := min(int(math.Ceil(1*float64(numClusters))), N)

	d := make(kmeans.Observations, N)
	for i := 0; i < N; i++ {
		row := make([]float64, numDims)
		for j, ind := range inds[:numDims] {
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

	clusterModesMap := map[int][]*mbc3.Mode{}
	modeClusterMap := map[*mbc3.Mode]int{}
	minRepeatedModes := N

	for i := 0; i < 1000; i++ {

		// Partition the data points
		clusters, err := km.Partition(d, numClusters, 0)
		if err != nil {
			return err
		}

		// Get cluster number for each mode (starts at 0)
		localClusterModesMap := map[int][]*mbc3.Mode{}
		localModeClusterMap := map[*mbc3.Mode]int{}
		for i, obs := range d {
			c := clusters.Nearest(obs)
			if _, ok := localClusterModesMap[c]; !ok {
				localClusterModesMap[c] = []*mbc3.Mode{}
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
		opModesMap := map[int][]*mbc3.Mode{}
		for _, m := range modes {
			opModesMap[m.OP] = append(opModesMap[m.OP], m)
		}

		// Rebuild list of modes keeping one for each OP
		modes = []*mbc3.Mode{}
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

//------------------------------------------------------------------------------
// Campbell Diagram
//------------------------------------------------------------------------------

type CampbellDiagram struct {
	HasWind    bool                  `json:"HasWind"`
	RotSpeeds  []float32             `json:"RotSpeeds"`
	WindSpeeds []float32             `json:"WindSpeeds"`
	Lines      []CampbellDiagramLine `json:"Lines"`
}

type CampbellDiagramLine struct {
	ID     int                    `json:"ID"`
	Label  string                 `json:"Label"`
	Points []CampbellDiagramPoint `json:"Points"`
	Hide   bool                   `json:"Hide"`
}

type CampbellDiagramPoint struct {
	OP            int     `json:"OpPtID"`
	Mode          int     `json:"ModeID"`
	RotSpeed      float32 `json:"RotSpeed"`
	WindSpeed     float32 `json:"WindSpeed"`
	NaturalFreqHz float32 `json:"NaturalFreqHz"`
	DampedFreqHz  float32 `json:"DampedFreqHz"`
	DampingRatio  float32 `json:"DampingRatio"`
}

// CampbellDiagram returns a Campbell Diagram structure from the results
func (res *Results) CampbellDiagram() *CampbellDiagram {

	hasWind := false
	rotSpeeds := make([]float32, len(res.OPs))
	windSpeeds := make([]float32, len(res.OPs))
	for i, op := range res.OPs {
		rotSpeeds[i] = float32(op.RotSpeed)
		windSpeeds[i] = float32(op.WindSpeed)
		hasWind = op.WindSpeed > 0 || hasWind
	}

	lines := make([]CampbellDiagramLine, len(res.ModeSets))
	for i, ms := range res.ModeSets {
		points := make([]CampbellDiagramPoint, len(ms.Modes))
		for j, m := range ms.Modes {
			points[j] = CampbellDiagramPoint{
				OP:            m.OP,
				Mode:          m.ID,
				RotSpeed:      float32(res.OPs[m.OP].RotSpeed),
				WindSpeed:     float32(res.OPs[m.OP].WindSpeed),
				NaturalFreqHz: float32(m.NaturalFreqHz),
				DampedFreqHz:  float32(m.DampedFreqHz),
				DampingRatio:  float32(m.DampingRatio),
			}
		}
		lines[i] = CampbellDiagramLine{
			ID:     i + 1,
			Label:  strconv.Itoa(i + 1),
			Points: points,
		}
	}

	return &CampbellDiagram{
		HasWind:    hasWind,
		RotSpeeds:  rotSpeeds,
		WindSpeeds: windSpeeds,
		Lines:      lines,
	}
}
