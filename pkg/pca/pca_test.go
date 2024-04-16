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
// Description: This file contains tests for the pca package.
package pca

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/mat"
)

// helper function to check if slices are almost equal
func slicesAlmostEqual(a, b []float64, tol float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if math.Abs(a[i]-b[i]) > tol {
			return false
		}
	}
	return true
}

func TestNIPALS(t *testing.T) {
	// Setup test data
	data := []float64{
		-1.18, -1.43, -1.17, -1.37, -1.61,
		-0.59, -0.99, -0.82, -1.12, -0.89,
		0.59, -0.44, -0.58, -0.93, -0.48,
		-1.18, 0.00, 0.23, 0.62, -0.16,
		0.00, 0.22, -0.35, 0.93, 0.89,
		0.59, 0.99, 0.70, 1.06, 1.05,
		1.77, 1.65, 1.99, 0.81, 1.21,
	}
	X := mat.NewDense(7, 5, data)

	// Perform PCA
	T, P, _, err := NIPALS(X, 5)
	if err != nil {
		t.Fatalf("NIPALS returned an error: %v", err)
	}

	// Expected values for the first principal component scores T
	expectedTScores := []float64{-3.034002, -1.980712, -0.874202, -0.167544, 0.763628, 1.977201, 3.315630}

	// Extract the first column of T (first principal component scores)
	actualTScores := mat.Col(nil, 0, T)

	// Expected values for the first principal component loadings P
	expectedPLoadings := []float64{0.390972, 0.486678, 0.454000, 0.426498, 0.471455}

	// Extract the first column of P (first principal component loadings)
	actualPLoadings := mat.Col(nil, 0, P)

	// Tolerance for floating-point comparison
	tolerance := 0.01

	// Compare actual and expected results for T and P
	if !slicesAlmostEqual(actualTScores, expectedTScores, tolerance) {
		t.Errorf("Principal component scores T do not match expected values. Got: %v, Want: %v", actualTScores, expectedTScores)
	}

	if !slicesAlmostEqual(actualPLoadings, expectedPLoadings, tolerance) {
		t.Errorf("Principal component loadings P do not match expected values. Got: %v, Want: %v", actualPLoadings, expectedPLoadings)
	}
}
