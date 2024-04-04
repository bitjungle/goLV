package pls

import (
	"errors"
	"fmt"
	"math"
	"math/rand"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

// InitializeScores initializes the score vector u for Y
func InitializeScores(Y *mat.Dense) *mat.Dense {
	rows, _ := Y.Dims()
	scores := make([]float64, rows)
	for i := range scores {
		scores[i] = rand.Float64() // Random initialization
	}
	return mat.NewDense(rows, 1, scores)
}

// NipalsPLS performs the NIPALS Algorithm for PLS Regression
func NipalsPLS(X, Y *mat.Dense, ncomp int, maxIter int, tol float64) (map[string]*mat.Dense, error) {
	rowsX, colsX := X.Dims()
	_, colsY := Y.Dims()
	//T, P, Q, W := make([]*mat.Dense, 0, ncomp), make([]*mat.Dense, 0, ncomp), make([]*mat.Dense, 0, ncomp), make([]*mat.Dense, 0, ncomp)
	P, Q, W := make([]*mat.Dense, 0, ncomp), make([]*mat.Dense, 0, ncomp), make([]*mat.Dense, 0, ncomp)
	T := mat.NewDense(rowsX, ncomp, nil) // Scores matrix

	for c := 0; c < ncomp; c++ {
		u := InitializeScores(Y) // Make sure this generates a column vector

		t := mat.NewDense(rowsX, 1, nil)
		p := mat.NewDense(colsX, 1, nil)
		q := mat.NewDense(colsY, 1, nil)
		u = mat.NewDense(rowsX, 1, nil)
		var tOld *mat.Dense
		for iteration := 0; iteration < maxIter; iteration++ {
			// Calculate w, normalize it
			w := mat.NewDense(colsX, 1, nil)
			w.Mul(X.T(), u)
			normalize(w)

			// Calculate t, normalize it, and check for convergence
			t.Mul(X, w)
			normalize(t)

			if iteration > 0 && normDiff(t, tOld) < tol {
				break
			}
			tOld = t

			// Calculate p, q
			p.Mul(X.T(), t)
			q.Mul(Y.T(), u)

			// Update u
			u.Mul(Y, q)

			// Store results
			//T = append(T, t)
			P = append(P, p)
			Q = append(Q, q)
			W = append(W, w)

		}
		T.SetCol(c, t.RawMatrix().Data) // Store the score vector in the scores matrix
		// Deflate X, Y
		deflate(X, t, p)
		deflate(Y, t, q)
	}

	// Tmat, err := StackDenseMatrices(T)
	// if err != nil {
	// 	return nil, err
	// }
	Pmat, err := StackDenseMatrices(P)
	if err != nil {
		return nil, err
	}
	Qmat, err := StackDenseMatrices(Q)
	if err != nil {
		return nil, err
	}
	Wmat, err := StackDenseMatrices(W)
	if err != nil {
		return nil, err
	}

	return map[string]*mat.Dense{
		"T": T,
		"P": Pmat,
		"Q": Qmat,
		"W": Wmat,
	}, nil
}

// PlsPredict makes predictions using a fitted NIPALS PLS model
func PlsPredict(XNew *mat.Dense, plsModel map[string]*mat.Dense) *mat.Dense {
	W, Q := plsModel["W"], plsModel["Q"]

	// Get dimensions
	rowsXNew, colsXNew := XNew.Dims()
	rowsW, colsW := W.Dims()
	rowsQ, colsQ := Q.Dims()

	// Log dimensions for debugging
	fmt.Printf("Multiplying TNew: XNew dimensions %d x %d, W dimensions %d x %d\n", rowsXNew, colsXNew, rowsW, colsW)

	// Projecting the new data onto the PLS components
	TNew := mat.NewDense(rowsXNew, colsW, nil)
	TNew.Mul(XNew, W)

	// Log dimensions for the second multiplication
	fmt.Printf("Multiplying YPred: TNew dimensions %d x %d, Q.T() dimensions %d x %d\n", rowsXNew, colsW, colsQ, rowsQ)

	// Making predictions using the loadings for Y
	YPred := mat.NewDense(rowsXNew, colsQ, nil)
	YPred.Mul(TNew, Q.T())

	return YPred
}

// StackDenseMatrices vertically stacks a slice of *mat.Dense matrices.
// All matrices must have the same number of columns.
func StackDenseMatrices(matrices []*mat.Dense) (*mat.Dense, error) {
	if len(matrices) == 0 {
		return nil, errors.New("no matrices to stack")
	}

	_, cols := matrices[0].Dims()
	var totalRows int
	for _, m := range matrices {
		rows, c := m.Dims()
		if c != cols {
			return nil, errors.New("matrices have different number of columns")
		}
		totalRows += rows
	}

	stacked := mat.NewDense(totalRows, cols, nil)
	currentRow := 0
	for _, m := range matrices {
		r, _ := m.Dims()
		stacked.Slice(currentRow, currentRow+r, 0, cols).(*mat.Dense).Copy(m)
		currentRow += r
	}

	return stacked, nil
}

// normalize modifies the matrix x to have unit length.
func normalize(x *mat.Dense) {
	data := x.RawMatrix().Data
	norm := floats.Norm(data, 2)
	if norm != 0 {
		floats.Scale(1/norm, data)
	}
}

// normDiff calculates the Euclidean norm of the difference between a and b.
func normDiff(a, b *mat.Dense) float64 {
	if a == nil || b == nil {
		return math.Inf(1)
	}

	r, c := a.Dims()
	diff := mat.NewDense(r, c, nil)
	diff.Sub(a, b)
	return mat.Norm(diff, 2)
}

// deflate subtracts the outer product of t and p from X.
func deflate(X, t, p *mat.Dense) {
	rows, cols := X.Dims()
	outer := mat.NewDense(rows, cols, nil)
	outer.Mul(t, p.T())
	X.Sub(X, outer)
}
