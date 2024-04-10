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
// Description: This file contains tests for the readdata package.
package readdata

import (
	"reflect"
	"testing"
)

// TestProcessCSV tests the ProcessCSV function using main_test_data.csv.
func TestProcessCSV(t *testing.T) {
	testDataFile := "../../data/read_test_data.csv"

	got, err := ProcessCSV(testDataFile)
	if err != nil {
		t.Fatalf("ProcessCSV() error = %v, wantErr nil", err)
	}

	// Define the expected result based on the contents of main_test_data.csv
	want := ProcessedData{
		VariableNames: []string{"Variable 1", "Variable2", "Variable3"},
		ObjectNames:   []string{"Object 1", "Object2", "Object3", "Object4", "Object5", "Object6", "Object7", "Object8", "Object9", "Object10"},
		Data: [][]float64{
			{0.8446018268164345, 0.888984055890554, 0.5305244852912209},
			{0.42713559719734184, 0.6259527055341845, 0.8272448514283607},
			{0.6923390114923365, 0.5723993819206581, 0.5862756192889184},
			{0.4584634554367081, 0.3582871787324112, 0.5098330389696466},
			{0.7725076275142637, 0.961167991594181, 0.1849102253321092},
			{0.5520770806761762, 0.40613512496823245, 0.05325803559269715},
			{1.0, 0.12200722164080213, 0.6804925451817727},
			{0.6283709864136304, 0.28029261420308016, 0.5755970412968474},
			{0.3998056116276213, 0.06814054967890837, 0.45520760119571435},
			{0.9714306431972592, 0.08139972989032118, 0.3625379743876117},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ProcessCSV() got = %v, want %v", got, want)
	}
}
