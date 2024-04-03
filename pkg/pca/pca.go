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
// Description: This file contains the NIPALS PCA algorithm and helper functions.
package pca

import (
	"math/rand"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

// NIPALS performs Principal Component Analysis (PCA) using the Non-linear Iterative Partial Least Squares (NIPALS) algorithm.
//
// X: Data matrix to perform PCA on.
// numComponents: Number of principal components to compute.
//
// Returns the scores matrix (T), loadings matrix (P), and the Eigenvalues.
//
// The NIPALS algorithm:
//
// Step 1: Choose a column vector t as an initial estimate of the first
// principal component score.
//
// Step 2: Compute the loading vector p as the matrix product of the
// transpose of t and X, and normalize p to have unit length.
//
// Step 3: Update the score vector t as the matrix product of X and p,
// and normalize t to have unit length.
//
// Step 4: Check the convergence of t and p by calculating the squared
// correlation between the old and new values. If the
// convergence criterion is met, stop the iteration. Otherwise,
// go back to Step 2.
//
// Step 5: Deflate the data matrix X by subtracting the outer product of
// t and p from X. This removes the variance captured by the
// first principal component.
//
// Step 6: Repeat Steps 1 to 5 to find the next principal component
// using the deflated X. Continue until the desired number of
// principal components is obtained.
func NIPALS(X mat.Matrix, numComponents int) (*mat.Dense, *mat.Dense, []float64, error) {
	epsilon := 1e-6
	maxIterations := 500

	rows, cols := X.Dims()

	T := mat.NewDense(rows, numComponents, nil)   // Scores matrix
	P := mat.NewDense(cols, numComponents, nil)   // Loadings matrix
	Eigenvalues := make([]float64, numComponents) // Eigenvalues for each component
	XRes := mat.DenseCopyOf(X)                    // Residual X matrix

	var t, p, tNew, outerProduct, sub mat.Dense

	for i := 0; i < numComponents; i++ { // Repeat for each component

		// Use the column from XRes with the highest variance as the initial t
		t.CloneFrom(initialScoreVector(XRes))
		// Use a random vector as the initial t
		//t.CloneFrom(initialRandomScoreVector(rows))

		for j := 0; j < maxIterations; j++ { // Repeat until convergence
			// Compute loading vector p
			p.Mul(XRes.T(), &t)

			// Normalize p to length 1
			pNorm := floats.Norm(p.RawMatrix().Data, 2)
			if pNorm == 0 {
				break // Avoid division by zero
			}
			p.Scale(1/pNorm, &p)

			// Compute score vector t
			tNew.Mul(XRes, &p)

			// Check for convergence
			if mat.Norm(&tNew, 2)-mat.Norm(&t, 2) < epsilon {
				break
			}
			t.CloneFrom(&tNew)
		}

		T.SetCol(i, t.RawMatrix().Data) // Store the score vector in the scores matrix
		P.SetCol(i, p.RawMatrix().Data) // Store the loading vector in the loadings matrix

		// Deflate the data matrix by the outer product of t and p
		outerProduct.Mul(&t, p.T())
		sub.Sub(XRes, &outerProduct)
		XRes.CloneFrom(&sub) // Update the residual X matrix

		// Calculate the eigenvalue for this principal component
		// Eigenvalue is approximated by the squared norm of the scores column
		tColumn := T.ColView(i)
		Eigenvalues[i] = mat.Norm(tColumn, 2)
		Eigenvalues[i] *= Eigenvalues[i]
	}

	return T, P, Eigenvalues, nil
}

// initialScoreVector selects the column from a given matrix that has the highest variance
func initialScoreVector(X *mat.Dense) *mat.Dense {
	rows, cols := X.Dims()
	maxVariance := 0.0
	var columnIndex int

	for j := 0; j < cols; j++ {
		var mean, variance float64
		for i := 0; i < rows; i++ {
			value := X.At(i, j)
			mean += value             // Accumulate the sum of values
			variance += value * value // Accumulate the sum of squares
		}
		mean /= float64(rows)                         // Calculate the mean
		variance = variance/float64(rows) - mean*mean // Calculate the variance

		if variance > maxVariance { // Check if this column has higher variance
			maxVariance = variance
			columnIndex = j
		}
	}

	// Extract the column with the highest variance
	highestVarianceColumn := mat.NewDense(rows, 1, nil)
	mat.Col(highestVarianceColumn.RawMatrix().Data, columnIndex, X)
	return highestVarianceColumn
}

// initialRandomScoreVector creates a random and normalized vector of scores
func initialRandomScoreVector(rows int) *mat.Dense {
	var t mat.Dense
	tRaw := make([]float64, rows) // Create a raw vector to store the random values
	for j := range tRaw {
		tRaw[j] = rand.Float64() // Assign a random value
	}
	tNorm := floats.Norm(tRaw, 2)   // Calculate the norm of the vector
	t = *mat.NewDense(rows, 1, nil) // Create a new matrix to store the normalized vector
	for j, v := range tRaw {
		t.Set(j, 0, v/tNorm) // Normalize the vector
	}
	return &t
}

// calculateVariancePercentages calculates the percentage of variance explained by each principal component.
func CalculateVariancePercentages(eigenvalues []float64) []float64 {
	sumEigenvalues := 0.0
	for _, eigenvalue := range eigenvalues {
		sumEigenvalues += eigenvalue
	}

	variancePercentages := make([]float64, len(eigenvalues))
	for i, eigenvalue := range eigenvalues {
		variancePercentages[i] = (eigenvalue / sumEigenvalues) * 100
	}

	return variancePercentages
}
