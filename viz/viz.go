package viz

import (
	"acdc/lin"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gonum.org/v1/gonum/cmplxs"
)

type Options struct {
	Scale float32
}

type ModeViz struct {
}

func (opts *Options) CalcViz(execPath string, rootPath string, mbc *lin.MBC, eigRes *lin.EigenResults) error {

	maxFreqHz := 10.0

	VTKLinTim := 2
	VTKLinTimes1 := true
	VTKLinPhase := 0.0

	// Collect modes at or below max frequency
	modes := []lin.Mode{}
	for _, m := range eigRes.Modes {
		if m.NaturalFreqHz <= maxFreqHz {
			modes = append(modes, m)
		}
	}

	// Collect mode data
	modeIDs := make([]string, len(modes))
	natFreq := make([]float64, len(modes))
	dampFreq := make([]float64, len(modes))
	dampRatio := make([]float64, len(modes))
	for i, m := range modes {
		modeIDs[i] = strconv.Itoa(m.ID + 1)
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
	for i, m := range modes {

		mags[i] = make([][]float64, len(mbc.Azimuths))
		phases[i] = make([][]float64, len(mbc.Azimuths))

		ev := make([]complex128, len(m.EigenVectorFull))

		// Loop through azimuths
		for j, azimuth := range mbc.Azimuths {

			for k := 0; k < numBlades; k++ {
				xi := azimuth + 2*math.Pi*float64(k)/float64(numBlades) // Blade angle
				s, c := math.Sincos(xi)
				tt[k] = [3]complex128{1, complex(c, 0), complex(s, 0)}
			}

			// Copy eigenvector
			copy(ev, m.EigenVectorFull)

			// If first value in eigenvector is negative, invert eigenvector
			// so all eigenvectors have the same sign (arbitrary)
			if real(ev[0]) < 0 {
				for i := range ev {
					ev[i] *= -1
				}
			}

			// Loop through all state triplets and convert non-rotating
			// eigenvectors back to rotating
			for _, triplet := range mbc.OrderX.Triplets {
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

			// Get magnitudes and phases of rotating eigenvector for this mode and azimuth
			mags[i][j] = make([]float64, len(ev))
			phases[i][j] = make([]float64, len(ev))
			for k, c := range ev {
				mags[i][j][k] = cmplx.Abs(c)
				phases[i][j][k] = cmplx.Phase(c)
			}
		}
	}

	vizFilePath := rootPath + ".ModeShapeVTK.viz"
	modesFilePath := rootPath + ".ModeShapeVTK.acdcMBC"
	modesFileName := filepath.Base(modesFilePath)
	chkpFileName := filepath.Base(rootPath) + ".ModeShapeVTK"

	// Write visualization file
	w := &bytes.Buffer{}
	fmt.Fprintf(w, "------- OpenFAST MODE-SHAPE INPUT FILE -------------------------------------------\n")
	fmt.Fprintf(w, "# Options for visualizing mode shapes\n")
	fmt.Fprintf(w, "---------------------- FILE NAMES ----------------------------------------------\n")
	fmt.Fprintf(w, "%-20s CheckpointRoot - Rootname of the checkpoint file written when OpenFAST generated the linearization files (without the \".chkp\" extension)\n", `"`+chkpFileName+`"`)
	fmt.Fprintf(w, "%-20s ModesFileName  - Name of the mode-shape file (with eigenvectors)\n", `"`+modesFileName+`"`)
	fmt.Fprintf(w, "---------------------- VISUALIZATION OPTIONS -----------------------------------\n")
	fmt.Fprintf(w, "%-20d VTKLinModes    - Number of modes to visualize (0 <= VTKLinModes <= NumModes)\n", numModes)
	fmt.Fprintf(w, "%-20s VTKModes       - List of which VTKLinModes modes will be visualized (modes will be added sequentially from the last value entered)\n", strings.Join(modeIDs, ", "))
	fmt.Fprintf(w, "%-20g VTKLinScale    - Mode shape visualization scaling factor (exaggerates mode shapes: try 10 for ElastoDyn; 0.1 for BeamDyn)\n", opts.Scale)
	fmt.Fprintf(w, "%-20d VTKLinTim      - Switch to make one animation for all LinTimes together (VTKLinTim=1) or separate animations for each LinTimes (VTKLinTim=2)\n", VTKLinTim)
	fmt.Fprintf(w, "%-20v VTKLinTimes1   - If VTKLinTim=2, visualize modes at LinTimes(1) only? (if false, files will be generated at all LinTimes)\n", VTKLinTimes1)
	fmt.Fprintf(w, "%-20f VTKLinPhase    - Phase used when making one animation for all LinTimes together (used only when VTKLinTim=1)\n", VTKLinPhase)
	if err := os.WriteFile(vizFilePath, w.Bytes(), 0777); err != nil {
		return err
	}

	// Write mode data binary file
	w.Reset()
	binary.Write(w, binary.LittleEndian, int32(1))
	binary.Write(w, binary.LittleEndian, int32(numModes))
	binary.Write(w, binary.LittleEndian, int32(len(mbc.DescStates)))
	binary.Write(w, binary.LittleEndian, int32(len(mbc.Azimuths)))
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
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Create command, get output pipe, set stderr to stdout, start command
	cmd := exec.CommandContext(ctx, execPath, "-VTKLin", vizFilePath)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
