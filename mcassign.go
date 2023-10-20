package main

import (
	"fmt"
	"math"
)

// Minimum Cost Assignment algorithm
// https://github.com/bmc/munkres
// https://brc2.com/the-algorithm-workshop/

const INVALID = math.MaxInt
const UNMARKED = 0
const STAR = 1
const PRIME = 2

type MCA struct {
	C          IntMatrix
	N          int
	Z0_r       int
	Z0_c       int
	Marked     IntMatrix
	path       [][2]int
	RowCovered []bool
	ColCovered []bool
	OrigSize   [2]int
}

func padMatrix(matrix [][]int, pad_value int) [][]int {

	// Calculate max number of rows and columns in matrix
	numRows := len(matrix)
	numCols := 0
	for _, row := range matrix {
		numCols = max(numCols, len(row))
	}

	// Get new square matrix size as max of number of columns or rows
	newSize := max(numCols, numRows)

	// Create new matrix and copy existing rows, pad row length if necessary
	newMatrix := NewIntMatrix(newSize, newSize, 0)
	for i, row := range matrix {
		copy(newMatrix[i], row)
		for j := numCols; j < newSize; j++ {
			newMatrix[i][j] = pad_value
		}
	}

	// If old matrix had fewer rows than columns, populate rows
	if numRows < newSize {
		newRow := make([]int, newSize)
		for i := range newRow {
			newRow[i] = pad_value
		}
		for i := numRows; i < newSize; i++ {
			copy(newMatrix[i], newRow)
		}
	}

	return newMatrix
}

type IntMatrix [][]int

func NewIntMatrix(m, n, v int) IntMatrix {
	rows := make([][]int, m)
	data := make([]int, n*m)
	for i := range data {
		data[i] = v
	}
	for i := range rows {
		rows[i], data = data[:n:n], data[n:]
	}
	return rows
}

// Compute the indexes for the lowest-cost pairings between rows and
// columns in the database. Returns a list of `(row, column)` tuples
// that can be used to traverse the matrix.

// - `cost_matrix` (list of lists of numbers): The cost matrix. If this
//   cost matrix is not square, it will be padded with zeros, via a call
//   to `padMatrix()`. (This method does *not* modify the caller's
//   matrix. It operates on a copy of the matrix.)

// A list of `(row, column)` tuples that describe the lowest cost path
// through the matrix
func MinCostAssignment(cost IntMatrix) (results [][2]int, err error) {

	// Pad cost matrix so it is square, get size
	costSq := padMatrix(cost, 0)
	N := len(costSq)

	// Create algorithm structure
	m := MCA{
		C:          costSq,
		OrigSize:   [2]int{len(cost), len(cost[0])},
		N:          N,
		RowCovered: make([]bool, N),
		ColCovered: make([]bool, N),
		Z0_r:       0,
		Z0_c:       0,
		Marked:     NewIntMatrix(N, N, 0),
		path:       make([][2]int, N),
	}

	done := false
	step := 1

	// Loop while not done
	for !done {

		// Switch based on step
		switch step {
		case 1:
			step, err = m.Step1()
		case 2:
			step, err = m.Step2()
		case 3:
			step, err = m.Step3()
		case 4:
			step, err = m.Step4()
		case 5:
			step, err = m.Step5()
		case 6:
			step, err = m.Step6()
		default:
			done = true
		}

		// If error occurred during step, return
		if err != nil {
			return nil, err
		}
	}

	// Assemble results from starred columns
	for i := 0; i < m.OrigSize[0]; i++ {
		for j := 0; j < m.OrigSize[1]; j++ {
			if m.Marked[i][j] == STAR {
				results = append(results, [2]int{i, j})
			}
		}
	}

	return results, nil
}

func (m IntMatrix) Copy() IntMatrix {
	mn := NewIntMatrix(len(m), len(m[0]), 0)
	for i, row := range m {
		copy(mn[i], row)
	}
	return mn
}

func (m *MCA) Step1() (int, error) {

	// Loop through rows in C
	for i := range m.C {

		// Find minimum value in row, ignore invalid values
		var minVal *int
		for j, v := range m.C[i] {
			if (v != INVALID) && (minVal == nil || v < *minVal) {
				minVal = &m.C[i][j]
			}
		}

		// If no min value found, return error
		if minVal == nil {
			return 0, fmt.Errorf("all values in row %d are INVALID", i+1)
		}

		// Subtract minimum value from all values in row
		for j, v := range m.C[i] {
			if v != INVALID {
				m.C[i][j] -= *minVal
			}
		}
	}

	// Go to Step 2
	return 2, nil
}

func (m *MCA) Step2() (int, error) {

	// Loop through rows of cost matrix
	for i := range m.C {

		// If row is covered (contains STARed 0), continue
		if m.RowCovered[i] {
			continue
		}

		// Loop through columns of cost matrix
		for j := range m.C[i] {

			// If column is covered (contains STARed 0), continue
			if m.ColCovered[j] {
				continue
			}

			// If element cost, mark as STAR and set row/column covered
			if m.C[i][j] == 0 {
				m.Marked[i][j] = STAR
				m.ColCovered[j] = true
				m.RowCovered[i] = true
				break
			}
		}
	}

	// Reset covered slices to false
	m.clearCovers()

	// Go to step 3
	return 3, nil
}

// Cover each column containing a starred zero. If K columns are
// covered, the starred zeros describe a complete set of unique
// assignments. In this case, Go to DONE, otherwise, Go to Step 4.
func (m *MCA) Step3() (int, error) {

	// Count each column containing a STARed 0
	count := 0

	// Loop through rows/columns of marked matrix
	for i := range m.Marked {
		for j := range m.Marked[i] {

			// If column is covered, continue
			if m.ColCovered[j] {
				continue
			}

			// If element is starred, mark column covered and increment count
			if m.Marked[i][j] == STAR {
				m.ColCovered[j] = true
				count += 1
			}
		}
	}

	// If all columns are covered, algorithm is done
	if count >= m.N {
		return 7, nil
	}

	// Go to step 4
	return 4, nil
}

// Find a non-covered zero and mark as PRIME. If there is no starred zero
// in the row containing this primed zero, Go to Step 5. Otherwise,
// cover this row and uncover the column containing the starred
// zero. Continue in this manner until there are no uncovered zeros
// left. Save the smallest uncovered value and Go to Step 6.
func (m *MCA) Step4() (int, error) {

	row := 0
	col := 0

	for {
		row, col := m.findZero(row, col)
		if row < 0 {
			return 6, nil
		}

		m.Marked[row][col] = PRIME

		star_col := m.findStarInRow(row)
		if star_col < 0 {
			m.Z0_r = row
			m.Z0_c = col
			return 5, nil
		}

		col = star_col
		m.RowCovered[row] = true
		m.ColCovered[col] = false
	}
}

// Construct a series of alternating primed and starred zeros as
// follows. Let Z0 represent the uncovered primed zero found in Step 4.
// Let Z1 denote the starred zero in the column of Z0 (if any).
// Let Z2 denote the primed zero in the row of Z1 (there will always
// be one). Continue until the series terminates at a primed zero
// that has no starred zero in its column. Un-star each starred zero
// of the series, star each primed zero of the series, erase all
// primes and uncover every line in the matrix. Return to Step 3
func (m *MCA) Step5() (int, error) {

	count := 0
	m.path[count][0] = m.Z0_r
	m.path[count][1] = m.Z0_c

	for {

		row := m.findStarInCol(m.path[count][1])
		if row < 0 {
			break
		}

		count += 1
		m.path[count][0] = row
		m.path[count][1] = m.path[count-1][1]

		col := m.findPrimeInRow(m.path[count][0])
		count += 1
		m.path[count][0] = m.path[count-1][0]
		m.path[count][1] = col
	}

	m.convertPath(m.path, count)
	m.clearCovers()
	m.erasePrimes()

	// Go to step 3
	return 3, nil
}

// Add the value found in Step 4 to every element of each covered
// row, and subtract it from every element of each uncovered column.
// Return to Step 4 without altering any stars, primes, or covered
// lines.
func (m *MCA) Step6() (int, error) {

	minVal := m.findSmallest()
	events := 0 // track actual changes to matrix
	for i := range m.C {
		for j := range m.C[i] {
			if m.C[i][j] == INVALID {
				continue
			}
			if m.RowCovered[i] {
				m.C[i][j] += minVal
				events += 1
			}
			if !m.ColCovered[j] {
				m.C[i][j] -= minVal
				events += 1
			}
			if m.RowCovered[i] && !m.ColCovered[j] {
				events -= 2 // change reversed, no real difference
			}
		}
	}
	if events == 0 {
		return 0, fmt.Errorf("matrix cannot be solved")
	}

	// Go to step 4
	return 4, nil
}

// Find the smallest uncovered value in the matrix.
func (m MCA) findSmallest() int {
	minVal := INVALID
	for i := range m.C {
		for j := range m.C[i] {
			if !m.RowCovered[i] && !m.ColCovered[j] {
				if m.C[i][j] != INVALID && minVal > m.C[i][j] {
					minVal = m.C[i][j]
				}
			}
		}
	}
	return minVal
}

// Find the first uncovered element with value 0
func (m MCA) findZero(i0, j0 int) (int, int) {

	row := -1
	col := -1
	i := i0
	done := false

	for !done {
		j := j0
		for {
			if (m.C[i][j] == 0) && (!m.RowCovered[i]) && (!m.ColCovered[j]) {
				row = i
				col = j
				done = true
			}
			j = (j + 1) % m.N
			if j == j0 {
				break
			}
		}
		i = (i + 1) % m.N
		if i == i0 {
			done = true
		}
	}

	return row, col
}

// Find the first starred element in the specified row. Returns
// the column index, or -1 if no starred element was found.
func (m MCA) findStarInRow(row int) int {
	for col, v := range m.Marked[row] {
		if v == STAR {
			return col
		}
	}
	return -1
}

// Find the first starred element in the specified row. Returns
// the row index, or -1 if no starred element was found.
func (m MCA) findStarInCol(col int) int {
	for row := range m.Marked {
		if m.Marked[row][col] == STAR {
			return row
		}
	}
	return -1
}

// Find the first prime element in the specified row. Returns
// the column index, or -1 if no starred element was found.
func (m MCA) findPrimeInRow(row int) int {
	for col := range m.C {
		if m.Marked[row][col] == PRIME {
			return col
		}
	}
	return -1
}

func (m *MCA) convertPath(path [][2]int, count int) {
	for i := 0; i < count+1; i++ {
		if m.Marked[path[i][0]][path[i][1]] == STAR {
			m.Marked[path[i][0]][path[i][1]] = UNMARKED
		} else {
			m.Marked[path[i][0]][path[i][1]] = STAR
		}
	}
}

// clearCovers resets row and column covered flags
func (m *MCA) clearCovers() {
	for i := range m.RowCovered {
		m.RowCovered[i] = false
	}
	for i := range m.ColCovered {
		m.ColCovered[i] = false
	}
}

// erasePrimes erases all prime markings
func (m *MCA) erasePrimes() {
	for i := range m.Marked {
		for j := range m.Marked[i] {
			if m.Marked[i][j] == PRIME {
				m.Marked[i][j] = UNMARKED
			}
		}
	}
}
