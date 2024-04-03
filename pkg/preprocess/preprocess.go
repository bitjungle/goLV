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
// Description: This file contains functions for preprocessing data for PCA/PLS.
package preprocess

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

// colMean calculates the mean of each column in a matrix.
func colMean(X *mat.Dense) []float64 {
	r, c := X.Dims()
	colMeans := make([]float64, c)
	for j := 0; j < c; j++ {
		col := X.ColView(j)
		colMeans[j] = mat.Sum(col) / float64(r)
	}
	return colMeans
}

// meanCenter centers the data by subtracting the mean of each column from its elements.
func MeanCenter(X *mat.Dense) (*mat.Dense, []float64) {
	r, c := X.Dims()       // Get the dimensions of the matrix
	colMeans := colMean(X) // Calling colMean internally

	centeredX := mat.NewDense(r, c, nil) // Create a new matrix to store the centered values
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			centeredX.Set(i, j, X.At(i, j)-colMeans[j])
		}
	}
	return centeredX, colMeans
}

// colStdDev calculates the standard deviation of each column in a matrix.
func colStdDev(X *mat.Dense) []float64 {
	r, c := X.Dims()
	colMeans := colMean(X)
	stdDevs := make([]float64, c)

	for j := 0; j < c; j++ {
		var sumSq float64   // Initialize sum of squares to zero
		col := X.ColView(j) // Get the column
		mean := colMeans[j] // Get the mean for the column

		for i := 0; i < r; i++ { // Loop over rows
			diff := col.AtVec(i) - mean
			sumSq += diff * diff
		}
		stdDevs[j] = math.Sqrt(sumSq / float64(r)) // This is the one to use
		//stdDevs[j] = math.Sqrt(sumSq / float64(r-1))
		//stdDevs[j] = math.Sqrt(sumSq / float64(r/(r-1)))
	}
	return stdDevs
}

// scaleByStdDev scales each column of the matrix by its standard deviation.
func ScaleByStdDev(X *mat.Dense) (*mat.Dense, []float64) {
	r, c := X.Dims()                   // Get the dimensions of the matrix
	colStd := colStdDev(X)             // Calculate standard deviations for each column
	scaledX := mat.NewDense(r, c, nil) // Create a new matrix to store the scaled values

	for j := 0; j < c; j++ {
		std := colStd[j]
		for i := 0; i < r; i++ {
			scaledVal := X.At(i, j) / std
			scaledX.Set(i, j, scaledVal)
		}
	}
	return scaledX, colStd
}

// autoscale centers the data by subtracting the mean of each column
// and then scales it by dividing by the standard deviation of each column.
func Autoscale(X *mat.Dense) (*mat.Dense, []float64, []float64) {
	centeredX, colMeans := MeanCenter(X)
	autoscaledX, colStd := ScaleByStdDev(centeredX)
	return autoscaledX, colMeans, colStd
}
