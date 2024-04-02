package utils

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

// SliceToDense converts a [][]float64 to a *mat.Dense matrix.
func SliceToDense(data [][]float64) *mat.Dense {
	if len(data) == 0 || len(data[0]) == 0 {
		return mat.NewDense(0, 0, nil)
	}

	rows := len(data)
	cols := len(data[0])

	// Flatten the [][]float64 into a []float64
	flatData := make([]float64, 0, rows*cols)
	for _, row := range data {
		if len(row) != cols {
			// Handle inconsistent row length
			fmt.Println("Error: Rows of unequal lengths")
			return nil
		}
		flatData = append(flatData, row...)
	}

	// Create a new mat.Dense matrix with the flattened data
	return mat.NewDense(rows, cols, flatData)
}

func PrettyPrintMatrix(matrix mat.Matrix, tit ...string) {
	title := "Matrix" // default title
	if len(tit) > 0 {
		title = tit[0] // if title is provided, use it
	}

	r, c := matrix.Dims()
	fmt.Printf("--- %s: Dimensions (%d, %d)\n", title, r, c)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			fmt.Printf("%9.6f ", matrix.At(i, j))
		}
		fmt.Println()
	}
	fmt.Println("---")
}
