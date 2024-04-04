package pls

import (
	"math/rand"
	"testing"

	"gonum.org/v1/gonum/mat"
)

// generateRandomData generates random data for testing.
func generateRandomData(rows, cols int) *mat.Dense {
	data := make([]float64, rows*cols)
	for i := range data {
		data[i] = rand.NormFloat64() // Random normal distribution
	}
	return mat.NewDense(rows, cols, data)
}

// TestNipalsPLS tests the NipalsPLS function of the pls package.
func TestNipalsPLS(t *testing.T) {
	// Define the dimensions of the data
	nSamples := 1000
	nFeatures := 100
	nTargets := 1
	nComponents := 5

	// Generate predictor variables (X)
	X := generateRandomData(nSamples, nFeatures)

	// Generate response variables (Y) using a simple linear relationship with the predictors
	Y := generateRandomData(nSamples, nTargets)

	// Perform PLS
	plsModel, err := NipalsPLS(X, Y, nComponents, 500, 1e-6)
	if err != nil {
		t.Fatalf("NipalsPLS failed: %v", err)
	}

	// Example usage with the same data for prediction (in a real test, use separate test data)
	YPred := PlsPredict(X, plsModel)

	// Basic validation of results
	if YPred == nil {
		t.Fatalf("Prediction failed, received nil matrix")
	}

	rows, cols := YPred.Dims()
	if rows != nSamples || cols != nTargets {
		t.Fatalf("Predicted matrix has incorrect dimensions: got (%d, %d), want (%d, %d)", rows, cols, nSamples, nTargets)
	}

	// Additional checks can be added here, like verifying the accuracy of the prediction
	// against known values or expected patterns in the data.
}
