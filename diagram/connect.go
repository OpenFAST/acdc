package diagram

import (
	"acdc/lin"
	"fmt"
	"math"
	"sort"

	"gonum.org/v1/gonum/mat"
)

// connectModesMAC builds connected sets of modes from linearization results
func connectModesMAC(OPs []lin.LinOP, freqRangeHz [2]float64, structMax bool) ([]*ModeSet, error) {

	// Create array of mode sets
	modeSets := []*ModeSet{}

	// Loop through modes in first operating point and add to map
	for i := range OPs[0].Modes {

		// Get pointer to mode
		m := &OPs[0].Modes[i]

		// If mode should not be filtered, continue
		if !m.Filter(freqRangeHz, structMax) {
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
		w := mat.NewDense(len(modeSets), len(op.Modes), nil)

		// Create map mapping mode index to mode
		modeIndexMap := map[int]*lin.Mode{}

		// Loop through modes in mode set map
		for j, modeSet := range modeSets {

			// 	Get last mode in mode set
			mp := modeSet.Modes[len(modeSet.Modes)-1]

			// Loop through modes in current operating point
			k := 0
			for l := range op.Modes {

				// Get mode
				mn := &op.Modes[l]

				// If mode should not be filtered, continue
				if !mn.Filter(freqRangeHz, structMax) {
					continue
				}

				// Calculate MAC between modes
				mac, err := mp.MAC(mn)
				if err != nil {
					return nil, err
				}

				// Modify MAC by change in frequency
				mac *= 1 - math.Abs(mn.NaturalFreqHz-mp.NaturalFreqHz)/(freqRangeHz[1]-freqRangeHz[0])

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
