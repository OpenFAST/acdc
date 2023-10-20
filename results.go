package main

import (
	"acdc/mbc3"
	"fmt"
	"math"
	"runtime"
	"sort"
	"strings"

	"github.com/mkmik/argsort"
	"github.com/parallelo-ai/kmeans"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
)

type Results struct {
	HasWind  bool        `json:"HasWind"`
	OPs      []OPResults `json:"OPs"`
	ModeSets []ModeSet   `json:"ModeSets"`
	MBC      []*mbc3.MBC `json:"-"`
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
	ID            int         `json:"ID"`
	Label         string      `json:"Label"`
	Weight        float64     `json:"Weight"`
	FrequencyMean float64     `json:"FrequencyMean"`
	Indices       []ModeIndex `json:"Indices"`
}

type ModeIndex struct {
	OP     int     `json:"OP"`
	Mode   int     `json:"Mode"`
	Weight float64 `json:"Weight"`
}

func LoadResults(LinFiles []string) (res *Results, err error) {

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

	// Sort results by group name
	sort.Slice(linFileResults, func(i, j int) bool {
		return linFileResults[i].Name < linFileResults[j].Name
	})

	// Loop through results and add
	res = &Results{}
	for i, lfr := range linFileResults {

		// Store data in results
		res.MBC = append(res.MBC, lfr.MBC)
		res.OPs = append(res.OPs, OPResults{
			ID:        i + 1,
			RotSpeed:  lfr.MBC.RotSpeed,
			WindSpeed: lfr.MBC.WindSpeed,
			Modes:     lfr.Modes,
		})

		// If non-zero wind speed in MBC, set flag that results have wind
		if lfr.MBC.WindSpeed > 0 {
			res.HasWind = true
		}
	}

	// Identify modes
	err = res.IdentifyModes2()
	if err != nil {
		return nil, err
	}

	// Identify modes
	err = res.IdentifyModes3()
	if err != nil {
		return nil, err
	}

	return res, nil
}

type LinFileGroup struct {
	Name  string
	Files []string
}

type LinFileResult struct {
	Name  string
	MBC   *mbc3.MBC
	Modes []mbc3.Mode
	err   error
}

func linFileWorker(linFilesChan <-chan LinFileGroup, resultsChan chan<- LinFileResult) {
	var err error

	// Loop through linearization files sent on channel
	for linFileGroup := range linFilesChan {

		// Read linearization files
		linFileData := make([]*mbc3.LinData, len(linFileGroup.Files))
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
		modes, err := mbc.EigenAnalysis()
		if err != nil {
			resultsChan <- LinFileResult{err: err}
			return
		}

		// Send MBC and mode results
		resultsChan <- LinFileResult{Name: linFileGroup.Name, MBC: mbc, Modes: modes}
	}
}

func (r *Results) IdentifyModes2() error {

	// Create a map of mode sets
	modeSetMap := map[int]*ModeSet{}
	modeSetMapNext := map[int]*ModeSet{}

	// Initialize mode set map with all modes from first OP
	for i := range r.OPs[0].Modes {
		modeSetMap[i] = &ModeSet{
			ID:      i + 1,
			Label:   fmt.Sprintf("%d", i+1),
			Indices: []ModeIndex{{OP: 0, Mode: i, Weight: 1}},
		}
	}

	// Loop through operating point results
	for i, op := range r.OPs {

		// Skip first operating point
		if i == 0 {
			continue
		}

		// Create weighting matrix
		w := mat.NewDense(len(modeSetMap), len(op.Modes), nil)

		// Loop through modes in mode set map
		for j, modeSet := range modeSetMap {

			// Loop through modes in next operating point
			for k, mn := range op.Modes {

				// 	Get last index in mode set
				ind := modeSet.Indices[len(modeSet.Indices)-1]

				// Get previous mode
				mp := r.OPs[ind.OP].Modes[ind.Mode]

				// Calculate MAC between modes
				mac, err := mp.MACXP(&mn)
				if err != nil {
					return err
				}

				// Add MAC to weight matrix
				w.Set(j, k, mac)
			}
		}

		// Get min, max, and range of weights
		wMin, wMax := mat.Min(w), mat.Max(w)
		wRange := wMax - wMin

		// Create cost matrix from weights (rescale to maximize precision)
		cost := NewIntMatrix(len(modeSetMap), len(op.Modes), 0)
		for j := range cost {
			for k := range cost[j] {
				v := w.At(j, k)
				cost[j][k] = int(1e7 * (1 - (v-wMin)/wRange))
			}
		}

		// Save cost matrix in operating point
		r.OPs[i].Costs = cost

		// Find mode pairings that minimizes the total cost
		pairs, err := MinCostAssignment(cost)
		if err != nil {
			return err
		}

		// Set mode connections
		for _, p := range pairs {

			// Look up mode set from previous mode index
			modeSet := modeSetMap[p[0]]

			// Add next operating point and mode combination
			modeSet.Indices = append(modeSet.Indices, ModeIndex{OP: i, Mode: p[1], Weight: w.At(p[0], p[1])})

			// Update mode set in next set map
			modeSetMapNext[p[1]] = modeSet
		}

		// Set next map to current map, clear next map
		modeSetMap, modeSetMapNext = modeSetMapNext, map[int]*ModeSet{}
	}

	// Clear mode sets in results
	r.ModeSets = []ModeSet{}

	// Loop through mode sets
	for _, modeSet := range modeSetMap {

		// Calculate mean frequency
		for _, ind := range modeSet.Indices {
			modeSet.FrequencyMean += r.OPs[ind.OP].Modes[ind.Mode].NaturalFreqHz
		}
		modeSet.FrequencyMean /= float64(len(modeSet.Indices))

		// Append mode set to results
		r.ModeSets = append(r.ModeSets, *modeSet)
	}

	// Sort mode sets by mean frequency
	sort.Slice(r.ModeSets, func(i, j int) bool {
		return r.ModeSets[i].FrequencyMean < r.ModeSets[j].FrequencyMean
	})

	return nil
}

func (r *Results) IdentifyModes() error {

	// Create mode graph
	g := ModeGraph{
		ids:                   make(map[*mbc3.Mode]int64),
		WeightedDirectedGraph: simple.NewWeightedDirectedGraph(0, math.Inf(1)),
	}

	// Loop through operating point results
	for i := range r.OPs {

		// Get pointer to operating point
		op := &r.OPs[i]

		// Get standard deviation of natural frequencies
		natFreq := make([]float64, len(op.Modes))
		for j, m := range op.Modes {
			natFreq[j] = m.NaturalFreqHz
		}
		natFreqStdDev := stat.StdDev(natFreq, nil)

		// Loop through modes in operating point, add nodes
		for j := range op.Modes {

			// Get pointer to mode data
			mode := &op.Modes[j]

			// Create node
			n := g.WeightedDirectedGraph.NewNode()
			nid := n.ID()
			n = ModeNode{
				id:   nid,
				OP:   i + 1,
				Mode: mode,
			}

			// Add node to graph
			g.WeightedDirectedGraph.AddNode(n)
			g.ids[mode] = nid

			// If first OP, continue
			if i == 0 {
				continue
			}

			// Loop through nodes in previous OP and connect
			for k := range r.OPs[i-1].Modes {

				// Get pointer to mode from previous OP
				modePrev := &r.OPs[i-1].Modes[k]

				// If difference between current mode and previous mode natural frequency
				// is greater than twice the frequency standard deviation, continue
				if math.Abs(mode.NaturalFreqHz-modePrev.NaturalFreqHz) > natFreqStdDev {
					continue
				}

				// Get previous node from previous mode pointer
				pn := g.nodeFor(modePrev)
				if pn == nil {
					continue
				}

				// Calculate weight
				w, err := modePrev.MAC(*mode)
				if err != nil {
					return err
				}

				// Add weighted edge from previous node to current node
				g.SetWeightedEdge(simple.WeightedEdge{F: pn, T: n, W: 1 / w})
			}
		}
	}

	//--------------------------------------------------------------------------
	// Locate modes across operating points
	//--------------------------------------------------------------------------

	// Create slice of slices of mode indices to contain identified paths
	r.ModeSets = []ModeSet{}

	// Get ending operating point
	opEnd := &r.OPs[len(r.OPs)-1]

	// Loop through operating point results to get path starting OPs
	for i := range r.OPs {

		// Get starting operating point
		opStart := &r.OPs[i]

		// Loop and find minimum path
		for {

			// Generate minimum path tree, must be regenerated as nodes are removed from graph
			// pathTree := path.DijkstraAllPaths(g)

			// nidStart := int64(-1)
			// nidEnd := int64(-1)
			minWeight := -1.0
			graphPath := []graph.Node{}

			// Loop through modes in starting operating point
			for j := range opStart.Modes {

				// Get graph node corresponding to mode, if already in path (nil), continue
				ns := g.nodeFor(&opStart.Modes[j])
				if ns == nil {
					continue
				}

				// Get all shortest paths from start node
				pathTree := path.DijkstraFrom(ns, g)

				// Loop through modes at last operating point
				for k := range opEnd.Modes {

					// Get graph node ID corresponding to ending mode,
					// If mode is not in graph (already in path), continue
					ne := g.nodeFor(&opEnd.Modes[k])
					if ne == nil {
						continue
					}

					// Get weight between start and end nodes
					p, weight := pathTree.To(ne.ID())
					if weight < minWeight || minWeight == -1.0 {
						minWeight = weight
						graphPath = p
					}
				}
			}

			if minWeight == -1 {
				break
			}

			// Convert graph path to min path, remove path nodes from graph
			// since multiple paths can't use the same node
			minPath := []ModeIndex{}
			for _, n := range graphPath {
				minPath = append(minPath, n.(ModeNode).Index())
				g.RemoveNode(n.ID())
			}

			// Add path to results
			r.ModeSets = append(r.ModeSets, ModeSet{
				Label:   fmt.Sprintf("%d", len(r.ModeSets)+1),
				Weight:  minWeight,
				Indices: minPath,
			})

			fmt.Println(r.ModeSets[len(r.ModeSets)-1])
		}
	}

	return nil
}

type ModeGraph struct {
	ids map[*mbc3.Mode]int64
	*simple.WeightedDirectedGraph
}

func (g ModeGraph) nodeFor(mode *mbc3.Mode) graph.Node {
	id, ok := g.ids[mode]
	if !ok {
		return nil
	}
	return g.WeightedDirectedGraph.Node(id)
}

type ModeNode struct {
	id   int64
	OP   int
	Mode *mbc3.Mode
}

func (n ModeNode) ID() int64        { return n.id }
func (n ModeNode) String() string   { return "" }
func (n ModeNode) Index() ModeIndex { return ModeIndex{OP: n.OP, Mode: n.Mode.ID} }

func (r *Results) IdentifyModes3() error {

	// Collect all modes in all operating points
	modes := []*mbc3.Mode{}
	for i := range r.OPs {
		for j := range r.OPs[i].Modes {
			m := &r.OPs[i].Modes[j]
			if m.NaturalFreqHz < 0.5 {
				m.OP = i
				modes = append(modes, m)
			}
		}
	}

	N := len(modes)

	// Create weight matrix by comparing modes with MAC
	Wf := make([][]float64, N)
	W := mat.NewDense(N, N, nil)
	D := mat.NewDense(N, N, nil)
	Dinv := mat.NewDense(N, N, nil)
	Disr := mat.NewDense(N, N, nil)
	for i, m1 := range modes {
		Wf[i] = make([]float64, N)
		for j, m2 := range modes {
			if i != j {
				if mac, err := m1.MACXP(m2); err == nil {
					W.Set(i, j, mac)
					Wf[i][j] = mac
				}
			}
		}
		di := mat.Sum(W.RowView(i))
		D.Set(i, i, di)
		Dinv.Set(i, i, 1/di)
		Disr.Set(i, i, 1/math.Sqrt(di))
	}

	// Calculate Laplacian matrix (D - W)
	L := mat.NewDense(N, N, nil)
	L.Sub(D, W)

	// Calculate random walk laplacian
	Lrw := mat.NewDense(N, N, nil)
	Lrw.Mul(Dinv, L)

	// Calculate the symmetric laplacian
	Lsym := mat.NewDense(N, N, nil)
	Lsym.Mul(Disr, L)
	Lsym.Mul(Lsym, Disr)

	eig := mat.Eigen{}
	if ok := eig.Factorize(Lsym, mat.EigenRight); !ok {
		return fmt.Errorf("error computing eigenvalues")
	}
	eigenValues := eig.Values(nil)
	eigenVectors := &mat.CDense{}
	eig.VectorsTo(eigenVectors)

	inds := argsort.SortSlice(eigenValues, func(i, j int) bool {
		return real(eigenValues[i]) < real(eigenValues[j])
	})

	numComps := 6

	// data := make([][]float64, N)
	// for i := range data {
	// 	row := make([]float64, numComps)
	// 	for j, ind := range inds[:numComps] {
	// 		row[j] = real(eigenVectors.At(i, ind))
	// 	}
	// 	data[i] = row
	// }

	// c, err := clusters.KMeans(1000, 4, clusters.EuclideanDistance)
	// if err != nil {
	// 	return err
	// }

	// err = c.Learn(data)
	// if err != nil {
	// 	return err
	// }

	// for i, clusterNum := range c.Guesses() {
	// 	modes[i].Cluster = clusterNum + 1
	// }

	var d kmeans.Observations
	for i := 0; i < N; i++ {
		row := make([]float64, numComps)
		for j, ind := range inds[:numComps] {
			row[j] = real(eigenVectors.At(i, ind))
		}
		d = append(d, kmeans.Coordinates(row))
	}

	// Partition the data points into 16 clusters
	km := kmeans.New()
	clusters, err := km.Partition(d, 6, 0)
	if err != nil {
		return err
	}

	for i, obs := range d {
		modes[i].Cluster = clusters.Nearest(obs)
	}

	return nil
}
