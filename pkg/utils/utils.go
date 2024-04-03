// Copyright (C) 2024 BITJUNGLE Rune Mathisen
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Description: This file contains misc. utility functions.
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

// DenseToSlice converts a *mat.Dense matrix to a [][]float64.
func DenseToSlice(d *mat.Dense) [][]float64 {
	r, c := d.Dims()
	data := make([][]float64, r)
	for i := 0; i < r; i++ {
		data[i] = make([]float64, c)
		for j := 0; j < c; j++ {
			data[i][j] = d.At(i, j)
		}
	}
	return data
}

// CreateFilledSlice creates a slice of float64 of a specified length, filled with a given value.
func CreateFilledSlice(length int, value float64) ([]float64, error) {
	if length < 0 {
		return nil, fmt.Errorf("length cannot be negative")
	}

	slice := make([]float64, length)
	for i := range slice {
		slice[i] = value
	}
	return slice, nil
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
