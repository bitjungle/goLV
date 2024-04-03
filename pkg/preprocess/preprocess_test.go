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
// Description: This file contains tests for the preprocess functions.
package preprocess_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/bitjungle/goLV/pkg/preprocess"
	"github.com/bitjungle/goLV/pkg/utils"
	"gonum.org/v1/gonum/mat"
)

func getTestData(dataType string) *mat.Dense {
	dataRaw := []float64{
		50, 67, 90, 98, 120,
		55, 71, 93, 102, 129,
		65, 76, 95, 105, 134,
		50, 80, 102, 130, 138,
		60, 82, 97, 135, 151,
		65, 89, 106, 137, 153,
		75, 95, 117, 133, 155,
	}
	dataCentered := []float64{
		-10, -13, -10, -22, -20,
		-5, -9, -7, -18, -11,
		5, -4, -5, -15, -6,
		-10, 0, 2, 10, -2,
		0, 2, -3, 15, 11,
		5, 9, 6, 17, 13,
		15, 15, 17, 13, 15,
	}
	dataAutoscaled := []float64{
		-1.18, -1.43, -1.17, -1.37, -1.61,
		-0.59, -0.99, -0.82, -1.12, -0.89,
		0.59, -0.44, -0.58, -0.93, -0.48,
		-1.18, 0.00, 0.23, 0.62, -0.16,
		0.00, 0.22, -0.35, 0.93, 0.89,
		0.59, 0.99, 0.70, 1.06, 1.05,
		1.77, 1.65, 1.99, 0.81, 1.21,
	}

	var data []float64
	switch dataType {
	case "raw":
		data = dataRaw
	case "centered":
		data = dataCentered
	case "autoscaled":
		data = dataAutoscaled
	default:
		data = dataRaw // Default to raw if an unrecognized type is specified
	}

	return mat.NewDense(7, 5, data)
}

// TestMeanCenter checks if MeanCenter function centers a matrix correctly.
func TestMeanCenter(t *testing.T) {
	X := getTestData("raw")
	utils.PrettyPrintMatrix(X, "Raw data")

	centered := preprocess.MeanCenter(X)
	expected := getTestData("centered")

	if !reflect.DeepEqual(centered, expected) {
		t.Errorf("Mean centering result was incorrect, got: %v, want: %v.", centered, expected)
	}
}

// TestAutoscale checks if Autoscale function correctly combines centering and scaling.
func TestAutoscale(t *testing.T) {
	X := getTestData("raw")

	autoscaled := preprocess.Autoscale(X)
	expected := getTestData("autoscaled")

	tolerance := 0.01 // Equal to two decimal places
	if !almostEqual(autoscaled, expected, tolerance) {
		t.Errorf("Autoscaling was incorrect, got: %v, want: %v.", autoscaled, expected)
	}
}

// almostEqual checks if two matrices are equal within a specified tolerance.
func almostEqual(a, b *mat.Dense, tol float64) bool {
	ra, ca := a.Dims()
	rb, cb := b.Dims()

	if ra != rb || ca != cb {
		return false
	}

	for i := 0; i < ra; i++ {
		for j := 0; j < ca; j++ {
			if math.Abs(a.At(i, j)-b.At(i, j)) > tol {
				return false
			}
		}
	}
	return true
}
