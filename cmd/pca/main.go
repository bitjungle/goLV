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
	"fmt"
	"log"
	"os"

	"github.com/bitjungle/goLV/pkg/pca"
	"github.com/bitjungle/goLV/pkg/preprocess"
	"github.com/bitjungle/goLV/pkg/readdata"
	"github.com/bitjungle/goLV/pkg/utils"
	"github.com/spf13/cobra"
	"gonum.org/v1/gonum/mat"
)

// AppVersion will be set at compile time using -ldflags
var AppVersion string

// Command line flags.
var (
	autoScaleFlag     bool
	numComponentsFlag int
	outputFile        string
)

// Results struct to hold PCA analysis results.
type Results struct {
	VariableNames       []string    `json:"variable_names"`
	ObjectNames         []string    `json:"object_names"`
	NumComponents       int         `json:"num_components"`
	Scores              [][]float64 `json:"scores"`
	Loadings            [][]float64 `json:"loadings"`
	Eigenvalues         []float64   `json:"eigenvalues"`
	VariancePercentages []float64   `json:"variance_percentages"`
	XMean               []float64   `json:"x_mean"`
	XStd                []float64   `json:"x_std"`
}

// main function sets up and runs the Cobra command line application.
func main() {
	var rootCmd = &cobra.Command{
		Use:   "pca",
		Short: "goLV Principal Component Analysis (PCA)",
		Long: `goLV Principal Component Analysis (PCA) - Copyright (C) 2024 BITJUNGLE Rune Mathisen. 
		        This program is distributed under the Apache license version 2.0`,
		Run: runRootCommand,
	}

	// Configuration of persistent flags for Cobra.
	rootCmd.PersistentFlags().IntVarP(&numComponentsFlag, "comps", "c", -1, "Number of principal components to compute")
	rootCmd.PersistentFlags().BoolVarP(&autoScaleFlag, "scale", "s", false, "Apply autoscaling")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "", "Path to output results as a JSON file (optional)")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Command execution error: %v", err)
	}
}

// runRootCommand is the primary function executed by Cobra on run.
func runRootCommand(cmd *cobra.Command, args []string) {
	fmt.Println("goLV Principal Component Analysis (PCA) version", AppVersion, "running...")
	fmt.Println()

	if len(args) < 1 {
		log.Fatal("Please provide a CSV file")
	}

	doAnalysis(args[0])
}

// loadData reads and processes CSV data.
func loadData(filename string) (readdata.ProcessedData, *mat.Dense, error) {
	records, err := readdata.ProcessCSV(filename)
	if err != nil {
		return readdata.ProcessedData{}, nil, err
	}
	X := utils.SliceToDense(records.Data)
	return records, X, nil
}

// determineNumComponents determines the number of PCA components to compute.
func determineNumComponents(X *mat.Dense) int {
	if numComponentsFlag <= 0 {
		_, cols := X.Dims()
		return cols // Default to the number of columns
	}
	return numComponentsFlag
}

// doAnalysis orchestrates the PCA analysis.
func doAnalysis(filename string) {
	// Load data
	records, X, err := loadData(filename)
	if err != nil {
		log.Fatalf("Error loading data: %v", err)
	}

	// Determine the number of components
	numComponents := determineNumComponents(X)

	// Preprocess the data (mean centering and optionally autoscaling)
	var Xpre *mat.Dense
	var Xmean, Xstd []float64
	if autoScaleFlag {
		Xpre, Xmean, Xstd = preprocess.Autoscale(X)
	} else {
		Xpre, Xmean = preprocess.MeanCenter(X)
		Xstd = make([]float64, Xpre.RawMatrix().Cols)
		for i := range Xstd {
			Xstd[i] = 1.0
		}
	}

	// Perform PCA
	T, P, eigv, err := pca.NIPALS(Xpre, numComponents)
	if err != nil {
		log.Fatalf("Error performing NIPALS PCA: %v", err)
	}
	variancePercentages := pca.CalculateVariancePercentages(eigv)

	// Prepare and output the results
	results := prepareResults(records, numComponents, T, P, eigv, variancePercentages, Xmean, Xstd)
	outputResults(results)
}

// prepareResults organizes PCA results into a structured format.
func prepareResults(records readdata.ProcessedData, numComponents int,
	T, P *mat.Dense, eigv, variancePercentages, Xmean, Xstd []float64) Results {
	return Results{
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
}

// outputResults handles outputting the results either to console or file.
func outputResults(results Results) {
	if outputFile != "" {
		saveResultsToFile(results, outputFile)
	} else {
		printResults(results)
	}
}

// saveResultsToFile saves PCA results to a JSON file.
func saveResultsToFile(results Results, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer file.Close()

	json.NewEncoder(file).Encode(results)
	fmt.Printf("Results saved to %s\n", filename)
}

// printResults displays PCA results in the console.
func printResults(results Results) {
	fmt.Printf("Variable names:\n%v\n", results.VariableNames)
	fmt.Printf("Object names:\n%v\n", results.ObjectNames)
	fmt.Printf("Number of components: %v\n", results.NumComponents)
	utils.PrettyPrintSlice(results.Scores, "Scores (T)")
	utils.PrettyPrintSlice(results.Loadings, "Loadings (P)")
	fmt.Printf("Eigenvalues:\n%v\n", results.Eigenvalues)
	fmt.Printf("Variance percentages:\n%v\n", results.VariancePercentages)
	fmt.Printf("X mean:\n%v\n", results.XMean)
	fmt.Printf("X std:\n%v\n", results.XStd)
}
