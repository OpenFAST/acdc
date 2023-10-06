package mbc3

import (
	"bytes"
	"fmt"
	"math"
	"os"

	"gonum.org/v1/gonum/mat"
)

type MBC struct {
	DescStates   []string
	NumDOF2      int
	NumDOF1      int
	RotSpeed     float64
	WindSpeed    float64
	OrderX       OPOrder
	OrderX2      OPOrder
	OrderX2dot   OPOrder
	OrderX1      OPOrder
	OrderU       OPOrder
	OrderY       OPOrder
	AvgA         *mat.Dense
	AvgX         *mat.VecDense
	AvgXdot      *mat.VecDense
	A_NR         []*mat.Dense
	B_NR         []*mat.Dense
	C_NR         []*mat.Dense
	D_NR         []*mat.Dense
	EigenResults *EigenResults
}

func Analyze(linPaths []string) (*MBC, error) {

	var err error

	// Read linearization data from files
	linFileData := make([]*LinData, len(linPaths))
	for i, f := range linPaths {
		if linFileData[i], err = ReadLinFile(f); err != nil {
			return nil, err
		}
	}

	// Create matrix data from linearization file data
	matData := NewMatData(linFileData)

	// Perform multi-blade coordinate transform
	mbc, err := matData.MBC3()
	if err != nil {
		return nil, err
	}

	// Get indices of eigenvector rows to keep (exclude X2dot)
	rows := append(mbc.OrderX2.Indices, mbc.OrderX1.Indices...)

	// Perform Eigen analysis on average A matrix
	eigenAnalysis, err := EigenAnalysis(mbc.AvgA, rows)
	if err != nil {
		return nil, err
	}

	// Save eigen results
	mbc.EigenResults = eigenAnalysis

	return mbc, err
}

func (md *MatData) MBC3() (*MBC, error) {

	// Set the number of blades (code currently support 3 bladed rotors)
	numBlades := 3

	// Create MBC structure
	mbc := MBC{}

	//--------------------------------------------------------------------------
	// State ordering (q2, qdot2, q1)
	//--------------------------------------------------------------------------

	if len(md.OP_x) > 0 {

		// Get sorted version of operating points (q2, q2dot, q1)
		opx := md.OP_x.Sort()

		// Count number of first order operating points (q1)
		NumX1 := 0
		for _, op := range opx {
			if op.DerivOrder == 1 {
				NumX1++
			}
		}

		// Number of second order operating points (q2 + q2dot)
		NumX2 := len(opx) - NumX1

		// Get order each group of state values
		// This assumes that half of the X2 points are q2 and half are q2dot
		mbc.OrderX2 = NewOPOrder(opx[:NumX2/2], numBlades)         // q2
		mbc.OrderX2dot = NewOPOrder(opx[NumX2/2:NumX2], numBlades) // q2dot
		mbc.OrderX1 = NewOPOrder(opx[NumX2:], numBlades)           // q1

		// Combine operating point orders
		mbc.OrderX = CombineOPOrders(mbc.OrderX2, mbc.OrderX2dot, mbc.OrderX1)
	}

	//--------------------------------------------------------------------------
	// Input Ordering
	//--------------------------------------------------------------------------

	if len(md.OP_u) > 0 {
		mbc.OrderU = NewOPOrder(md.OP_u, numBlades)
	}

	//--------------------------------------------------------------------------
	// Output Ordering
	//--------------------------------------------------------------------------

	if len(md.OP_y) > 0 {
		mbc.OrderY = NewOPOrder(md.OP_y, numBlades)
	}

	//--------------------------------------------------------------------------
	// Blade triplet reordering
	//--------------------------------------------------------------------------

	// Create row permutation array for states
	P := mat.NewDense(mbc.OrderX.Num, mbc.OrderX.Num, nil)
	P.Permutation(mbc.OrderX.Num, mbc.OrderX.Indices)

	// Loop through linearization data
	for i := 0; i < md.NumStep; i++ {

		// Rotor speed in radians/sec and rotor speed squared
		omega := md.Omega[i]
		omegaDot := 0.0

		// Calculate t_tilde matrices
		tt := mat.NewDense(3, numBlades, nil)
		tt2 := mat.NewDense(3, numBlades, nil)
		tt3 := mat.NewDense(3, numBlades, nil)
		for j := 0; j < numBlades; j++ {
			xi := md.AzimuthRad[i] + 2*math.Pi*float64(j)/float64(numBlades) // Blade angle (Eq. 1)
			s, c := math.Sincos(xi)
			tt.SetRow(j, []float64{1, c, s})    // t_tilde   (Eq. 9)
			tt2.SetRow(j, []float64{0, -s, c})  // t_tilde_2 (Eq. 16a)
			tt3.SetRow(j, []float64{0, -c, -s}) // t_tilde_3 (Eq. 16b)
		}

		// t_tilde inverse
		ttv := &mat.Dense{}
		ttv.Inverse(tt)

		X2IFixed := eye(mbc.OrderX2.NumFixed)
		X20Fixed := &mat.Dense{}
		if mbc.OrderX2.NumFixed > 0 {
			X20Fixed.ReuseAs(mbc.OrderX2.NumFixed, mbc.OrderX2.NumFixed)
		}
		X1IFixed := eye(mbc.OrderX1.NumFixed)
		X10Fixed := &mat.Dense{}
		if mbc.OrderX1.NumFixed > 0 {
			X10Fixed.ReuseAs(mbc.OrderX1.NumFixed, mbc.OrderX1.NumFixed)
		}

		// Equation 11 for second-order states only
		T1 := blockDiag(X2IFixed, Repeat(tt, mbc.OrderX2.NumTriplets)...)

		// Inverse of T1
		T1v := blockDiag(X2IFixed, Repeat(ttv, mbc.OrderX2.NumTriplets)...)

		// Equation 14  for second-order states only
		T2 := blockDiag(X20Fixed, Repeat(tt2, mbc.OrderX2.NumTriplets)...)
		T2_omega := &mat.Dense{}
		T2_omegaDot := &mat.Dense{}
		T2_2omega := &mat.Dense{}
		if mbc.OrderX2.Num > 0 {
			T2_omega.Scale(omega, T2)
			T2_omegaDot.Scale(omegaDot, T2)
			T2_2omega.Scale(2*omega, T2)
		}

		// Equation 11 for first-order states (Equation 8 in MBC3 Update document)
		T1q := blockDiag(X1IFixed, Repeat(tt, mbc.OrderX1.NumTriplets)...)

		// Inverse of T1q
		T1qv := blockDiag(X1IFixed, Repeat(ttv, mbc.OrderX1.NumTriplets)...)

		// Equation 14 for first-order states (Equation  9 in MBC3 Update document)
		T2q_omega := &mat.Dense{}
		T2q := blockDiag(X10Fixed, Repeat(tt2, mbc.OrderX1.NumTriplets)...)
		if mbc.OrderX1.Num > 0 {
			T2q_omega.Scale(omega, T2q)
		}

		// Equation 15
		T3 := blockDiag(X20Fixed, Repeat(tt3, mbc.OrderX2.NumTriplets)...)
		T3_omega2 := &mat.Dense{}
		if mbc.OrderX2.Num > 0 {
			T3_omega2.Scale(omega*omega, T3)
		}

		// T1c := &mat.Dense{}
		// if mbc.OrderU.NumFixed > 0 {
		// 	T1c = blockDiag(eye(mbc.OrderU.NumFixed), Repeat(tt, mbc.OrderU.Num)...)
		// }

		// // Inverse of T1q
		// T1ov := &mat.Dense{}
		// if mbc.OrderY.NumFixed > 0 {
		// 	T1ov = blockDiag(eye(mbc.OrderY.NumFixed), Repeat(ttv, mbc.OrderY.Num)...)
		// }

		// Copy A matrix from linearization data
		A := mat.DenseCopyOf(md.A[i])
		A.Mul(P, A)     // Reorder rows
		A.Mul(A, P.T()) // Reorder columns

		// Equation 29 [[T1, 0, 0], [Omega*T2, T1, 0], [0, 0, T1q]]
		L := blockDiag(T1, T1, T1q)
		L.Slice(mbc.OrderX2.Num, mbc.OrderX2.Num+mbc.OrderX2dot.Num, 0, mbc.OrderX2.Num).(*mat.Dense).Copy(T2_omega)

		// Equation 29 [[omega*T2, 0, 0], [omega^2*T3+omegadot*T2, 2*omega*T2, 0], [0, 0, omega*T2q]]
		R := blockDiag(T2_omega, T2_2omega, T2q_omega)
		R.Slice(mbc.OrderX2.Num, mbc.OrderX2.Num+mbc.OrderX2dot.Num, 0, mbc.OrderX2.Num).(*mat.Dense).Add(T3_omega2, T2_omegaDot)

		// Equation 29
		ANR := &mat.Dense{}
		ANR.Mul(A, L)
		ANR.Sub(ANR, R)
		ANR.Mul(blockDiag(T1v, T1v, T1qv), ANR)
		ANR.Mul(P.T(), ANR) // Restore row order
		ANR.Mul(ANR, P)     // Restore column order

		// Save non-rotating A matrix
		mbc.A_NR = append(mbc.A_NR, ANR)
	}

	// Average the A matrix
	mbc.AvgA = mat.NewDense(len(md.OP_x), len(md.OP_x), nil)
	for _, A_NR := range mbc.A_NR {
		mbc.AvgA.Add(mbc.AvgA, A_NR)
	}
	mbc.AvgA.Scale(1/float64(md.NumStep), mbc.AvgA)

	// Average X operating points
	mbc.AvgX = mat.NewVecDense(len(md.OP_x), nil)
	for i := 0; i < md.NumStep; i++ {
		mbc.AvgX.AddVec(mbc.AvgX, md.OpX.RowView(i))
	}
	mbc.AvgX.ScaleVec(1/float64(md.NumStep), mbc.AvgX)

	// Average Xdot operating points
	mbc.AvgXdot = mat.NewVecDense(len(md.OP_x), nil)
	for i := 0; i < md.NumStep; i++ {
		mbc.AvgXdot.AddVec(mbc.AvgXdot, md.OpXd.RowView(i))
	}
	mbc.AvgXdot.ScaleVec(1/float64(md.NumStep), mbc.AvgXdot)

	return &mbc, nil
}

func eye(n int) *mat.Dense {
	if n == 0 {
		return &mat.Dense{}
	}
	d := make([]float64, n*n)
	for i := 0; i < n*n; i += n + 1 {
		d[i] = 1
	}
	return mat.NewDense(n, n, d)
}

func blockDiag(base *mat.Dense, other ...*mat.Dense) *mat.Dense {

	mats := append([]*mat.Dense{base}, other...)

	size := 0
	for _, m := range mats {
		_, c := m.Dims()
		size += c
	}

	if size == 0 {
		return &mat.Dense{}
	}

	M := mat.NewDense(size, size, nil)
	c := 0
	for _, m := range mats {
		_, cm := m.Dims()
		if cm > 0 {
			M.Slice(c, c+cm, c, c+cm).(*mat.Dense).Copy(m)
			c += cm
		}
	}

	return M
}

func Repeat[T any](item T, n int) []T {
	s := make([]T, n)
	for i := range s {
		s[i] = item
	}
	return s
}

func ToCSV(m mat.Matrix, path, format string) error {
	buf := &bytes.Buffer{}
	r, c := m.Dims()
	// fmt.Fprintf(buf, "%d,%d\n", r, c)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			if j > 0 {
				buf.WriteString(",")
			}
			fmt.Fprintf(buf, format, m.At(i, j))
		}
		buf.WriteString("\n")
	}
	return os.WriteFile(path, buf.Bytes(), 0777)
}
