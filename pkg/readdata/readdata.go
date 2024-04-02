// Description: This package processing provides functions to process CSV data
// for PCA. It includes functionalities to read CSV files, and convert data to
// float64.
package readdata

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ProcessedData encapsulates the variable names, object names, and converted
// float data from a CSV file.
type ProcessedData struct {
	VariableNames []string    // Variable names from the first row
	ObjectNames   []string    // Object names from the first column
	Data          [][]float64 // Data converted to float64
}

// ReadCSV reads a CSV file and returns a 2D slice of strings representing the
// data. The first row is assumed to be headers, and the first column is assumed
// to contain object names.
func ReadCSV(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	return reader.ReadAll()
}

// ProcessCSV reads data from a CSV file and returns variable names, object
// names, and the data as floats. The first row is assumed to contain variable
// names, and the first column in each row is assumed to contain object names.
func ProcessCSV(filename string) (ProcessedData, error) {
	records, err := ReadCSV(filename)
	if err != nil {
		return ProcessedData{}, err
	}

	// Check for sufficient data
	if len(records) < 2 || len(records[0]) < 2 {
		return ProcessedData{}, fmt.Errorf("CSV file must contain at least one row and one column of data")
	}

	variableNames := records[0][1:] // Skip the first cell
	var objectNames []string
	var floatData [][]float64 // Corrected type to [][]float64

	for _, record := range records[1:] { // Skip the first row (header)
		objectNames = append(objectNames, record[0])
		floatRow, err := convertToFloats(record[1:]) // Skip the first column (object name)
		if err != nil {
			return ProcessedData{}, err
		}
		floatData = append(floatData, floatRow) // Append floatRow correctly
	}

	return ProcessedData{
		VariableNames: variableNames,
		ObjectNames:   objectNames,
		Data:          floatData,
	}, nil
}

// convertToFloats converts a slice of strings to a slice of float64.
// An error is returned if any string cannot be converted to a float.
func convertToFloats(strs []string) ([]float64, error) {
	var floats []float64
	for _, str := range strs {
		trimmedStr := strings.TrimSpace(str) // Trim spaces from the string
		f, err := strconv.ParseFloat(trimmedStr, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing float in convertToFloats: %v", err)
		}
		floats = append(floats, f)
	}
	return floats, nil
}
