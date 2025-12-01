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
	"sort"
	"strconv"
	"strings"

	"github.com/dominikbraun/graph"
	"gonum.org/v1/gonum/cmplxs"
)

type Options struct {
	Scale float32
}

type Point struct {
	XYZ [3]float32 `json:"XYZ"`
}

type Component struct {
	Line      []Point `json:"Line"`
	LocalLine []Point `json:"LocalLine"`
}

type Frame struct {
	Components map[string]*Component `json:"Components"`
}

type ModeData struct {
	LineID int     `json:"LineID"`
	OPID   int     `json:"OPID"`
	ModeID int     `json:"ModeID"`
	Frames []Frame `json:"Frames"`
}

func (opts *Options) GenerateModeData(execPath string, op *lin.LinOP, modeIDs []int) (*ModeData, error) {

	// Get slice of modes from mode IDs
	modes := []lin.Mode{}
	for _, mID := range modeIDs {
		modes = append(modes, op.Modes[mID])
	}

	// Collect mode data
	natFreq := make([]float64, len(modes))
	dampFreq := make([]float64, len(modes))
	dampRatio := make([]float64, len(modes))
	for i, m := range modes {
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

	// Create file names and paths
	checkpointFileName := filepath.Base(op.RootPath) + ".ModeShapeVTK"
	vizFilePath := op.RootPath + ".ModeShapeVTK.viz"
	modesFilePath := op.RootPath + ".ModeShapeVTK.acdcMBC"
	modesFileName := filepath.Base(modesFilePath)

	// Write visualization file
	w := &bytes.Buffer{}
	fmt.Fprintf(w, "------- OpenFAST MODE-SHAPE INPUT FILE -------------------------------------------\n")
	fmt.Fprintf(w, "# Options for visualizing mode shapes\n")
	fmt.Fprintf(w, "---------------------- FILE NAMES ----------------------------------------------\n")
	fmt.Fprintf(w, "%-20v CheckpointRoot - Rootname of the checkpoint file written when OpenFAST generated the linearization files (without the \".chkp\" extension)\n", `"`+checkpointFileName+`"`)
	fmt.Fprintf(w, "%-20v ModesFileName  - Name of the mode-shape file (with eigenvectors)\n", `"`+modesFileName+`"`)
	fmt.Fprintf(w, "---------------------- VISUALIZATION OPTIONS -----------------------------------\n")
	fmt.Fprintf(w, "%-20v VTKLinModes    - Number of modes to visualize (0 <= VTKLinModes <= NumModes)\n", numModes)
	fmt.Fprintf(w, "%-20v VTKModes       - List of which VTKLinModes modes will be visualized (modes will be added sequentially from the last value entered)\n", 1)
	fmt.Fprintf(w, "%-20v VTKLinScale    - Mode shape visualization scaling factor (exaggerates mode shapes: try 10 for ElastoDyn; 0.1 for BeamDyn)\n", opts.Scale)
	fmt.Fprintf(w, "%-20v VTKLinTim      - Switch to make one animation for all LinTimes together (VTKLinTim=1) or separate animations for each LinTimes (VTKLinTim=2)\n", 2)
	fmt.Fprintf(w, "%-20v VTKLinTimes1   - If VTKLinTim=2, visualize modes at LinTimes(1) only? (if false, files will be generated at all LinTimes)\n", true)
	fmt.Fprintf(w, "%-20v VTKLinPhase    - Phase used when making one animation for all LinTimes together (used only when VTKLinTim=1)\n", 0.0)
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
	binary.Write(w, binary.LittleEndian, dampRatio)
	binary.Write(w, binary.LittleEndian, dampFreq)
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

	// Open log file for writing
	logFile, err := os.Create(op.RootPath + ".ModeShapeVTK.log")
	if err != nil {
		return nil, err
	}
	defer logFile.Close()

	// Get the case directory
	projectDir := filepath.Dir(filepath.Dir(filepath.Dir(vizFilePath)))

	// Get relative path from project directory to main file
	relPath, err := filepath.Rel(projectDir, vizFilePath)
	if err != nil {
		return nil, err
	}

	// Create command to generate vtk output and run it
	cmd := exec.Command(execPath, "-VTKLin", relPath)
	cmd.Dir = projectDir
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	// Get list of VTP files for this mode visualization
	rootName := filepath.Base(op.RootPath)
	vtpFilePaths, err := filepath.Glob(filepath.Join(filepath.Dir(op.RootPath), "vtk", rootName+"*.vtp"))
	if err != nil {
		return nil, err
	}

	// Parse mode data from files
	md, err := ParseModeData(vtpFilePaths)

	if err != nil {
		return nil, err
	}

	return md, nil
}

func ParseModeData(vtpFilePaths []string) (*ModeData, error) {

	sort.Strings(vtpFilePaths)

	// Create mode viz struct
	mv := ModeData{}

	// Loop through files
	for _, vtpFile := range vtpFilePaths {
		fmt.Println("\nProcessing file:", vtpFile)

		// Skip BD blade rotating states files
		if strings.Contains(filepath.Base(vtpFile), "BD_BldMotionRot") {
			continue
		}

		// Load vtk file
		vtk, local_vtk, err := LoadVTK(vtpFile)
		if err != nil {
			return nil, err
		}

		// Skip files without lines
		// TODO: add handling files only containing points
		if vtk.PolyData.Piece.NumberOfLines == 0 {
			continue
		}

		// Split file name
		tmp := strings.Split(filepath.Base(vtpFile), ".")

		// Get frame number for file
		frameNum, err := strconv.Atoi(tmp[len(tmp)-2])
		if err != nil {
			return nil, err
		}

		// If mode viz has fewer frames than frame number, append empty frames
		if len(mv.Frames) < frameNum {
			newFrames := make([]Frame, frameNum-len(mv.Frames))
			for i := range newFrames {
				newFrames[i].Components = make(map[string]*Component)
			}
			mv.Frames = append(mv.Frames, newFrames...)
		}

		// Get pointer to frame corresponding to this file
		frame := &mv.Frames[frameNum-1]

		// Get component name from file name
		componentName := tmp[len(tmp)-3]

		// Get component, if it doesn't exist create it
		component, ok := frame.Components[componentName]
		if !ok {
			component = &Component{}
			frame.Components[componentName] = component
		}

		// Generate list of connectivity
		connectivity := [][]int32{}
		offsetStart := 0
		for _, offsetEnd := range vtk.PolyData.Piece.Lines.DataArray[1].Offsets {
			connectivity = append(connectivity, vtk.PolyData.Piece.Lines.DataArray[0].Connectivity[offsetStart:offsetEnd])
			offsetStart = int(offsetEnd)
		}

		// Build graph and get sorted connectivity
		g := graph.New(graph.IntHash, graph.Directed())
		for _, c := range vtk.PolyData.Piece.Lines.DataArray[0].Connectivity {
			g.AddVertex(int(c))
		}
		for _, conn := range connectivity {
			for i := 0; i < len(conn)-1; i++ {
				g.AddEdge(int(conn[i]), int(conn[i+1]))
			}
		}
		conn, err := graph.TopologicalSort(g)
		if err != nil {
			return nil, err
		}

		// Copy line data into component
		component.Line = make([]Point, len(conn))
		for j, c := range conn {
			copy(component.Line[j].XYZ[:], vtk.PolyData.Piece.Points.DataArray.MatrixF32[c])
		}

		// Copy local line data into component
		if local_vtk != nil {
			component.LocalLine = make([]Point, len(conn))
			for j, c := range conn {
				copy(component.LocalLine[j].XYZ[:], local_vtk.PolyData.Piece.Points.DataArray.MatrixF32[c])
			}
		}
	}

	return &mv, nil
}
