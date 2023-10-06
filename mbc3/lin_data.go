package mbc3

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/mat"
)

type LinData struct {
	ID         int
	SimTime    float64 // sec
	RotorSpeed float64 // rad/s
	Azimuth    float64 // rad
	WindSpeed  float64 // m/s

	Num_x    int
	Num_xdot int
	Num_x2   int
	Num_xd   int
	Num_z    int
	Num_u    int
	Num_y    int

	OP_x, OP_xdot OPSlice
	OP_u, OP_y    OPSlice
	OP_xd, OP_z   OPSlice

	A, B, C, D *mat.Dense
	dUdu, dUdy *mat.Dense
}

func NewLinData(ID int) *LinData {
	return &LinData{
		ID:         ID,
		RotorSpeed: math.NaN(),
		WindSpeed:  math.NaN(),
	}
}

// OPData contains the operating point data.
type OPData struct {
	Index      int
	Value      float64
	IsRotFrame bool
	DerivOrder int
	Desc       string
}

// SortKey returns the operating point key for sorting
func (op OPData) SortKey() int {
	key := 1000
	switch op.DerivOrder {
	case 2:
		key = 10
		// If description starts with First time derivative
		_, desc, _ := strings.Cut(op.Desc, " ")
		if strings.HasPrefix(desc, "First time derivative") {
			key += 5
		}
	case 1:
		key = 20
	case 0:
		key = 30
	}
	// Rotating frames after non-rotating
	if op.IsRotFrame {
		key += 1
	}
	return key
}

// SanitizedDesc returns the description with text between parentheses removed
func (op OPData) SanitizedDesc() string {
	if strings.HasPrefix(op.Desc, "ED") {
		if j := strings.Index(op.Desc, "("); j > -1 {
			if k := strings.LastIndex(op.Desc, ")"); k > -1 {
				return op.Desc[:j] + op.Desc[k+1:]
			}
		}
	}
	return op.Desc
}

type OPSlice []OPData

func (ops OPSlice) Values() []float64 {
	vs := make([]float64, len(ops))
	for i, op := range ops {
		vs[i] = op.Value
	}
	return vs
}

func (ops OPSlice) DerivOrders() []int {
	vs := make([]int, len(ops))
	for i, op := range ops {
		vs[i] = op.DerivOrder
	}
	return vs
}

func (ops OPSlice) Indices() []int {
	vs := make([]int, len(ops))
	for i, op := range ops {
		vs[i] = op.Index
	}
	return vs
}

func (ops OPSlice) Descs() []string {
	vs := make([]string, len(ops))
	for i, op := range ops {
		vs[i] = op.Desc
	}
	return vs
}

func (ops OPSlice) Sort() OPSlice {

	// Create copy of operating points
	opsSorted := append(make(OPSlice, 0, len(ops)), ops...)

	// Sort operating points into q2, q2dot, and q1,
	// original order is maintained if categories are equal
	sort.SliceStable(opsSorted, func(i, j int) bool {
		return opsSorted[i].SortKey() < opsSorted[j].SortKey()
	})

	return opsSorted
}

// ReadLinFile reads the linearization file and parses the data. File name
// must end in ".N.lin" where N is a number
func ReadLinFile(filePath string) (*LinData, error) {

	// Open linearization file
	linFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer linFile.Close()

	// Get linearization file ID
	tmp := strings.Split(filepath.Base(filePath), ".")
	if len(tmp) < 3 {
		return nil, fmt.Errorf("invalid lin file name, must end in '.N.lin'")
	}
	ID, err := strconv.Atoi(tmp[len(tmp)-2])
	if err != nil {
		return nil, fmt.Errorf("invalid lin file name, must end in '.N.lin'")
	}

	// Create scanner to read linearization file
	scanner := bufio.NewScanner(linFile)

	// Initialize linearization data structure
	ld := &LinData{ID: ID}

	//--------------------------------------------------------------------------
	// Header
	//--------------------------------------------------------------------------

	for scanner.Scan() {

		// Get line without leading/trailing whitespace, skip if empty
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		// Split line into fields
		fields := strings.Fields(line)

		if strings.HasPrefix(line, "Simulation time") {
			if ld.SimTime, err = strconv.ParseFloat(fields[2], 64); err != nil {
				return nil, fmt.Errorf("error parsing Simulation time: %w", err)
			}
		} else if strings.HasPrefix(line, "Rotor Speed") {
			if ld.RotorSpeed, err = strconv.ParseFloat(fields[2], 64); err != nil {
				return nil, fmt.Errorf("error parsing Rotor Speed: %w", err)
			}
		} else if strings.HasPrefix(line, "Azimuth") {
			if ld.Azimuth, err = strconv.ParseFloat(fields[1], 64); err != nil {
				return nil, fmt.Errorf("error parsing Azimuth: %w", err)
			}
			ld.Azimuth = math.Mod(ld.Azimuth, 2*math.Pi)
		} else if strings.HasPrefix(line, "Wind Speed") {
			if ld.WindSpeed, err = strconv.ParseFloat(fields[2], 64); err != nil {
				return nil, fmt.Errorf("error parsing Wind Speed: %w", err)
			}
		} else if strings.HasPrefix(line, "Number of continuous states") {
			if ld.Num_x, err = strconv.Atoi(fields[4]); err != nil {
				return nil, fmt.Errorf("error parsing Number of continuous states: %w", err)
			}
		} else if strings.HasPrefix(line, "Number of discrete states") {
			if ld.Num_xd, err = strconv.Atoi(fields[4]); err != nil {
				return nil, fmt.Errorf("error parsing Number of discrete states: %w", err)
			}
		} else if strings.HasPrefix(line, "Number of constraint states") {
			if ld.Num_z, err = strconv.Atoi(fields[4]); err != nil {
				return nil, fmt.Errorf("error parsing Number of constraint states: %w", err)
			}
		} else if strings.HasPrefix(line, "Number of inputs") {
			if ld.Num_u, err = strconv.Atoi(fields[3]); err != nil {
				return nil, fmt.Errorf("error parsing Number of inputs: %w", err)
			}
		} else if strings.HasPrefix(line, "Number of outputs") {
			if ld.Num_y, err = strconv.Atoi(fields[3]); err != nil {
				return nil, fmt.Errorf("error parsing Number of outputs: %w", err)
			}
		} else if strings.HasPrefix(line, "Jacobians included") {
			break
		}
	}

	//--------------------------------------------------------------------------
	// Operating points
	//--------------------------------------------------------------------------

	var currentOP *OPSlice
	hasDeriv := false
	defaultDeriv := 0

	for scanner.Scan() {

		line := strings.TrimSpace(scanner.Text())

		switch line {
		case "":
			continue
		case "Order of continuous states:":
			currentOP = &ld.OP_x
			defaultDeriv = 2
			continue
		case "Order of continuous state derivatives:":
			currentOP = &ld.OP_xdot
			defaultDeriv = 2
			continue
		case "Order of discrete states:":
			currentOP = &ld.OP_x
			defaultDeriv = 2
			continue
		case "Order of inputs:":
			currentOP = &ld.OP_u
			defaultDeriv = 0
			continue
		case "Order of outputs:":
			currentOP = &ld.OP_y
			defaultDeriv = 0
			continue
		}

		if strings.Contains(line, "Operating Point") {
			hasDeriv = strings.Contains(line, "Derivative Order")
			continue
		}

		if line == "Linearized state matrices:" || line == "Jacobian matrices:" {
			break
		}

		fields := strings.Fields(line)
		field := ""

		// Create operating point structure
		op := OPData{DerivOrder: defaultDeriv}

		// Get first column as integer, skip line if not valid
		field, fields = fields[0], fields[1:]
		if op.Index, err = strconv.Atoi(field); err != nil {
			continue
		}
		// Decrement index to be zero based
		op.Index--

		// Get operating point value
		field, fields = fields[0], fields[1:]
		if strings.Contains(field, ",") {
			field = strings.Trim(field, ",")
			fields = fields[2:]
		}
		if op.Value, err = strconv.ParseFloat(field, 64); err != nil {
			return nil, err
		}

		// Get if OP is in the rotating frame
		field, fields = fields[0], fields[1:]
		if op.IsRotFrame, err = strconv.ParseBool(field); err != nil {
			return nil, err
		}

		// If OP has derivative order
		if hasDeriv {
			field, fields = fields[0], fields[1:]
			if op.DerivOrder, err = strconv.Atoi(field); err != nil {
				return nil, err
			}
		}

		// Combine remaining fields into description
		op.Desc = strings.Join(fields, " ")

		// Add operating piont to current OP
		*currentOP = append(*currentOP, op)

	}

	// Sum number of second order continuous states
	for _, op := range ld.OP_x {
		if op.DerivOrder == 2 {
			ld.Num_x2++
		}
	}

	//--------------------------------------------------------------------------
	// Matrices
	//--------------------------------------------------------------------------

	var matrix *mat.Dense
	iRow := 0
	for scanner.Scan() {

		// Get line with whitespace removed, skip empty line
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		// Skip linearized line
		if line == "Linearized state matrices:" {
			continue
		}

		// Split line by whitespace
		fields := strings.Fields(line)

		// If this line specifies the name of a matrix and its name
		if len(fields) == 4 && fields[2] == "x" {
			rows, _ := strconv.Atoi(fields[1])
			cols, _ := strconv.Atoi(fields[3])
			matrix = mat.NewDense(rows, cols, nil)
			iRow = 0

			switch fields[0] {
			case "A:":
				ld.A = matrix
			case "B:":
				ld.B = matrix
			case "C:":
				ld.C = matrix
			case "D:":
				ld.D = matrix
			case "dUdu:":
				ld.dUdu = matrix
			case "dUdy:":
				ld.dUdy = matrix
			}

			continue
		}

		row := make([]float64, len(fields))
		for i, s := range fields {
			if row[i], err = strconv.ParseFloat(s, 64); err != nil {
				return nil, err
			}
		}
		matrix.SetRow(iRow, row)
		iRow++
	}

	// Get error from scanner
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ld, nil
}
