package viz

import (
	"acdc/lin"
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/cmplxs"
)

type Options struct {
	Scale float32
}

type ModeViz struct {
}

func (opts *Options) CalcViz(execPath string, op *lin.OPResult, modeIDs []int) (*ModeViz, error) {

	// Create checkpoint file name
	checkpointFileName := filepath.Base(op.RootPath) + ".ModeShapeVTK"

	// Get slice of modes from mode IDs
	modes := []lin.Mode{}
	for _, mID := range modeIDs {
		modes = append(modes, op.EigRes.Modes[mID])
	}

	// Collect mode data
	modeStr := make([]string, len(modes))
	natFreq := make([]float64, len(modes))
	dampFreq := make([]float64, len(modes))
	dampRatio := make([]float64, len(modes))
	for i, m := range modes {
		modeStr[i] = strconv.Itoa(m.ID + 1)
		natFreq[i] = m.NaturalFreqHz
		dampFreq[i] = m.DampedFreqHz
		dampRatio[i] = m.DampingRatio
	}
	numModes := len(modes)

	const numBlades = 3
	vnr := [3]complex128{}
	vr := [3]complex128{}
	tt := [3][numBlades]complex128{}
	mags := make([][][]float64, len(modes))
	phases := make([][][]float64, len(modes))

	// Loop through modes
	for i, m := range modes {

		// Initialize magnitudes and phases for this mode
		mags[i] = make([][]float64, len(op.MBC.Azimuths))
		phases[i] = make([][]float64, len(op.MBC.Azimuths))

		ev := make([]complex128, len(m.EigenVector))

		// Loop through azimuths
		for j, azimuth := range op.MBC.Azimuths {

			// Construct tt matrix for converting from non-rotating eigenvectors to rotating
			for k := 0; k < numBlades; k++ {
				xi := azimuth + 2*math.Pi*float64(k)/float64(numBlades) // Blade angle
				s, c := math.Sincos(xi)
				tt[k] = [3]complex128{1, complex(c, 0), complex(s, 0)}
			}

			// Copy eigenvector from mode for modification
			copy(ev, m.EigenVector)

			// If first value in eigenvector is negative, invert eigenvector
			// so all eigenvectors have the same sign (arbitrary)
			if real(ev[0]) < 0 {
				for i := range ev {
					ev[i] *= -1
				}
			}

			// Loop through all state triplets and convert non-rotating
			// eigenvectors back to rotating
			for _, triplet := range op.MBC.OrderX.Triplets {
				for k, ind := range triplet {
					vnr[k] = ev[ind]
				}
				for k := range vnr {
					vr[k] = cmplxs.Dot(tt[k][:], vnr[:])
				}
				for k, ind := range triplet {
					ev[ind] = vr[k]
				}
			}

			// Initialize magnitude and phase array for this eigenvector
			mags[i][j] = make([]float64, len(ev))
			phases[i][j] = make([]float64, len(ev))

			// Get magnitudes and phases of rotating eigenvector for this mode and azimuth
			for k, c := range ev {
				mags[i][j][k] = cmplx.Abs(c)
				phases[i][j][k] = cmplx.Phase(c)
			}
		}
	}

	vizFilePath := op.RootPath + ".ModeShapeVTK.viz"
	modesFilePath := op.RootPath + ".ModeShapeVTK.acdcMBC"
	modesFileName := filepath.Base(modesFilePath)

	// Write visualization file
	w := &bytes.Buffer{}
	fmt.Fprintf(w, "------- OpenFAST MODE-SHAPE INPUT FILE -------------------------------------------\n")
	fmt.Fprintf(w, "# Options for visualizing mode shapes\n")
	fmt.Fprintf(w, "---------------------- FILE NAMES ----------------------------------------------\n")
	fmt.Fprintf(w, "%-20s CheckpointRoot - Rootname of the checkpoint file written when OpenFAST generated the linearization files (without the \".chkp\" extension)\n", `"`+checkpointFileName+`"`)
	fmt.Fprintf(w, "%-20s ModesFileName  - Name of the mode-shape file (with eigenvectors)\n", `"`+modesFileName+`"`)
	fmt.Fprintf(w, "---------------------- VISUALIZATION OPTIONS -----------------------------------\n")
	fmt.Fprintf(w, "%-20d VTKLinModes    - Number of modes to visualize (0 <= VTKLinModes <= NumModes)\n", numModes)
	fmt.Fprintf(w, "%-20s VTKModes       - List of which VTKLinModes modes will be visualized (modes will be added sequentially from the last value entered)\n", strings.Join(modeStr, ", "))
	fmt.Fprintf(w, "%-20g VTKLinScale    - Mode shape visualization scaling factor (exaggerates mode shapes: try 10 for ElastoDyn; 0.1 for BeamDyn)\n", opts.Scale)
	fmt.Fprintf(w, "%-20d VTKLinTim      - Switch to make one animation for all LinTimes together (VTKLinTim=1) or separate animations for each LinTimes (VTKLinTim=2)\n", 2)
	fmt.Fprintf(w, "%-20v VTKLinTimes1   - If VTKLinTim=2, visualize modes at LinTimes(1) only? (if false, files will be generated at all LinTimes)\n", true)
	fmt.Fprintf(w, "%-20f VTKLinPhase    - Phase used when making one animation for all LinTimes together (used only when VTKLinTim=1)\n", 0.0)
	if err := os.WriteFile(vizFilePath, w.Bytes(), 0777); err != nil {
		return nil, err
	}

	// Write mode data binary file
	w.Reset()
	binary.Write(w, binary.LittleEndian, int32(1))
	binary.Write(w, binary.LittleEndian, int32(numModes))
	binary.Write(w, binary.LittleEndian, int32(len(op.MBC.DescStates)))
	binary.Write(w, binary.LittleEndian, int32(len(op.MBC.Azimuths)))
	binary.Write(w, binary.LittleEndian, natFreq)
	binary.Write(w, binary.LittleEndian, dampFreq)
	binary.Write(w, binary.LittleEndian, dampRatio)
	for i := range modes {
		for _, azEVMags := range mags[i] {
			binary.Write(w, binary.LittleEndian, azEVMags)
		}
		for _, azEVPhases := range phases[i] {
			binary.Write(w, binary.LittleEndian, azEVPhases)
		}
	}
	if err := os.WriteFile(modesFilePath, w.Bytes(), 0777); err != nil {
		return nil, err
	}

	// Create command to generate vtk output and run it
	cmd := exec.Command(execPath, "-VTKLin", vizFilePath)
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return &ModeViz{}, nil
}
