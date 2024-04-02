package main

import (
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
	fmt.Println("Copyright (C) 2004 - BITJUNGLE Rune Mathisen")
	fmt.Println("This program is distributed under the Apache license version 2.0")
	fmt.Println()

	numComponentsFlag := flag.Int("comps", -1, "Number of principal components to compute")
	autoScaleFlag := flag.Bool("scale", false, "Apply autoscaling")
	versionFlag := flag.Bool("v", false, "Prints the version information")

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

	records, err := readdata.ProcessCSV(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Print input data
	log.Println("Variable names:", records.VariableNames)
	log.Println("Object names:", records.ObjectNames)
	X := utils.SliceToDense(records.Data)
	utils.PrettyPrintMatrix(X, "Data (X)")

	// Determine the number of components based on the flag or default to number of columns
	numComponents := *numComponentsFlag
	if numComponents <= 0 {
		_, numComponents = X.Dims() // Default to the number of columns
	}

	var Xpre *mat.Dense
	if *autoScaleFlag {
		Xpre = preprocess.Autoscale(X)
	} else {
		Xpre = preprocess.MeanCenter(X)
	}

	// Perform NIPALS PCA
	T, P, eigv, err := pca.NIPALS(Xpre, numComponents)
	if err != nil {
		log.Fatalf("Error performing NIPALS PCA: %v", err)
	}

	// Display results
	utils.PrettyPrintMatrix(T, "Scores (T)")
	utils.PrettyPrintMatrix(P, "Loadings (P)")
	fmt.Printf("Eigenvalues:\n%v\n", eigv)
	variancePercentages := pca.CalculateVariancePercentages(eigv)
	fmt.Printf("Variance percentages:\n%v\n", variancePercentages)

}
