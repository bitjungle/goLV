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
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/bitjungle/goLV/pkg/pca"
	"github.com/bitjungle/goLV/pkg/preprocess"
	"github.com/bitjungle/goLV/pkg/readdata"
	"github.com/bitjungle/goLV/pkg/utils"
	"gonum.org/v1/gonum/mat"
)

// AppVersion will be set at compile time using -ldflags
var AppVersion string

func main() {
	// Print copyright and license information
	fmt.Println("goLV Principal Component Analysis (PCA) version ", AppVersion)
	fmt.Println("Copyright (C) 2024 - BITJUNGLE Rune Mathisen")
	fmt.Println("This program is distributed under the Apache license version 2.0")
	fmt.Println()

	numComponentsFlag := flag.Int("comps", -1, "Number of principal components to compute")
	autoScaleFlag := flag.Bool("scale", false, "Apply autoscaling")
	versionFlag := flag.Bool("v", false, "Prints the version information")
	outputFile := flag.String("output", "", "Path to output results as a JSON file (optional)")

	flag.Parse()

	// If the version flag is present, print the version and exit
	if *versionFlag {
		fmt.Printf("golv version: %s\n", AppVersion)
		os.Exit(0)
	}

	if len(flag.Args()) < 1 {
		log.Fatal("Please provide a CSV file")
	}

	filename := flag.Arg(0)

	// Read and prepare input data
	records, err := readdata.ProcessCSV(filename)
	if err != nil {
		log.Fatal(err)
	}
	X := utils.SliceToDense(records.Data)

	// Determine the number of components based on the flag or default to number of columns
	numComponents := *numComponentsFlag
	if numComponents <= 0 {
		_, numComponents = X.Dims() // Default to the number of columns
	}

	var Xpre *mat.Dense
	var Xmean []float64
	var Xstd []float64
	if *autoScaleFlag {
		Xpre, Xmean, Xstd = preprocess.Autoscale(X)
	} else {
		Xpre, Xmean = preprocess.MeanCenter(X)
		Xstd, _ = utils.CreateFilledSlice(Xpre.RawMatrix().Cols, 1.0)
	}

	// Perform NIPALS PCA
	T, P, eigv, err := pca.NIPALS(Xpre, numComponents)
	if err != nil {
		log.Fatalf("Error performing NIPALS PCA: %v", err)
	}
	variancePercentages := pca.CalculateVariancePercentages(eigv)

	// Prepare results
	results := struct {
		VariableNames       []string    `json:"variable_names"`
		ObjectNames         []string    `json:"object_names"`
		NumComponents       int         `json:"num_components"`
		Scores              [][]float64 `json:"scores"`
		Loadings            [][]float64 `json:"loadings"`
		Eigenvalues         []float64   `json:"eigenvalues"`
		VariancePercentages []float64   `json:"variance_percentages"`
		XMean               []float64   `json:"x_mean"`
		XStd                []float64   `json:"x_std"`
	}{
		VariableNames:       records.VariableNames,
		ObjectNames:         records.ObjectNames,
		NumComponents:       numComponents,
		Scores:              utils.DenseToSlice(T),
		Loadings:            utils.DenseToSlice(P),
		Eigenvalues:         eigv,
		VariancePercentages: variancePercentages,
		XMean:               Xmean,
		XStd:                Xstd,
	}

	// Output results to JSON file if outputFile is provided, otherwise print to console
	if *outputFile != "" {
		file, err := os.Create(*outputFile)
		if err != nil {
			log.Fatalf("Failed to create output file: %v", err)
		}
		defer file.Close()
		json.NewEncoder(file).Encode(results)
		fmt.Printf("Results saved to %s\n", *outputFile)
	} else {
		// Display results
		fmt.Printf("Variable names:\n%v\n", results.VariableNames)
		fmt.Printf("Object names:\n%v\n", results.ObjectNames)
		fmt.Printf("Number of components: %v\n", results.NumComponents)
		utils.PrettyPrintMatrix(T, "Scores (T)")
		utils.PrettyPrintMatrix(P, "Loadings (P)")
		fmt.Printf("Eigenvalues:\n%v\n", results.Eigenvalues)
		fmt.Printf("Variance percentages:\n%v\n", results.VariancePercentages)
		fmt.Printf("X mean:\n%v\n", results.XMean)
		fmt.Printf("X std:\n%v\n", results.XStd)
	}
}
