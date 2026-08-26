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

		// Collect the modes in this operating point that pass the filter.
		// NOTE: these are gathered up front so the weight matrix has exactly
		// one column per candidate. Sizing it from len(op.Modes) leaves
		// trailing all-zero columns that are never written but are still
		// visible to mat.Max below.
		filteredModes := []*lin.Mode{}
		for l := range op.Modes {
			mn := &op.Modes[l]
			if mn.Filter(freqRangeHz, structMax) {
				filteredModes = append(filteredModes, mn)
			}
		}

		// No candidate modes in this operating point, nothing to connect
		if len(filteredModes) == 0 {
			continue
		}

		// No mode sets to connect to yet, because no mode in any earlier
		// operating point passed the filter. Seed one set per candidate mode;
		// building a zero-row weight matrix below would panic.
		if len(modeSets) == 0 {
			for _, mn := range filteredModes {
				modeSets = append(modeSets, &ModeSet{
					ID:    len(modeSets),
					Label: fmt.Sprintf("%d", len(modeSets)),
					Modes: []*lin.Mode{mn},
				})
			}
			continue
		}

		// Create empty weighting matrix
		w := mat.NewDense(len(modeSets), len(filteredModes), nil)

		// Loop through mode sets
		for j, modeSet := range modeSets {

			// 	Get last mode in mode set
			mp := modeSet.Modes[len(modeSet.Modes)-1]

			// Loop through candidate modes in current operating point
			for k, mn := range filteredModes {

				// Calculate MAC between modes
				mac, err := mp.MAC(mn)
				if err != nil {
					return nil, err
				}

				// Modify MAC by change in frequency
				mac *= 1 - math.Abs(mn.NaturalFreqHz-mp.NaturalFreqHz)/(freqRangeHz[1]-freqRangeHz[0])

				// Add MAC to weight matrix
				w.Set(j, k, mac)
			}
		}

		// Get max weight value
		wMax := mat.Max(w)

		// Create cost matrix (ints) from weights (rescale to maximize
		// precision). If nothing correlates at all then every cost is equal,
		// and the guard is required because dividing by a zero wMax yields
		// NaN, whose conversion to int is not defined by the language spec.
		cost := NewIntMatrix(len(modeSets), len(filteredModes), 0)
		if wMax > 0 {
			for j := range cost {
				for k := range cost[j] {
					cost[j][k] = int(1e7 * (1 - w.At(j, k)/wMax))
				}
			}
		}

		// Find mode pairings that minimizes the total cost
		pairs, err := MinCostAssignment(cost)
		if err != nil {
			return nil, err
		}

		// Add connected modes to sets, tracking which candidates were paired
		paired := make([]bool, len(filteredModes))
		for _, pair := range pairs {

			// Look up mode set from previous mode index
			modeSet := modeSets[pair[0]]

			// Add paired mode to slice of modes
			modeSet.Modes = append(modeSet.Modes, filteredModes[pair[1]])

			// Mark paired candidate mode
			paired[pair[1]] = true
		}

		// Loop through unpaired modes and create new mode sets.
		// NOTE: walk the slice in index order rather than ranging over a map.
		// Go randomizes map iteration order, so the IDs and labels assigned to
		// these new mode sets varied between runs on identical input.
		for k, mn := range filteredModes {
			if paired[k] {
				continue
			}
			modeSets = append(modeSets, &ModeSet{
				ID:    len(modeSets),
				Label: fmt.Sprintf("%d", len(modeSets)),
				Modes: []*lin.Mode{mn},
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
