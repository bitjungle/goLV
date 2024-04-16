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
package utils_test

import (
	"reflect"
	"testing"

	"github.com/bitjungle/goLV/pkg/utils"
)

// TestCreateFilledSlice tests the CreateFilledSlice function.
func TestCreateFilledSlice(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name    string
		length  int
		value   float64
		want    []float64
		wantErr bool
	}{
		{"Zero Length", 0, 5.0, []float64{}, false},
		{"Positive Length", 3, 2.5, []float64{2.5, 2.5, 2.5}, false},
		{"Negative Length", -1, 3.0, nil, true},
	}

	for _, tc := range testCases {
		got, err := utils.CreateFilledSlice(tc.length, tc.value)
		if (err != nil) != tc.wantErr {
			t.Errorf("CreateFilledSlice() error = %v, wantErr %v", err, tc.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("CreateFilledSlice() got = %v, want %v", got, tc.want)
		}
	}
}

func TestNormalize(t *testing.T) {
	data := [][]float64{
		{-1.18, -1.43, -1.17, -1.37, -1.61},
		{-0.59, -0.99, -0.82, -1.12, -0.89},
		{0.59, -0.44, -0.58, -0.93, -0.48},
		{-1.18, 0., 0.23, 0.62, -0.16},
		{0., 0.22, -0.35, 0.93, 0.89},
		{0.59, 0.99, 0.7, 1.06, 1.05},
		{1.77, 1.65, 1.99, 0.81, 1.21},
	}
	utils.PrettyPrintSlice(data)
	// dataD := utils.SliceToDense(data)
	// got := utils.Normalize(dataD)
	// want := [][]float64{
	// 	{-0.4472136, -0.54166667, -0.44211739, -0.51688178, -0.60857062},
	// 	{-0.2236068, -0.375, -0.30986005, -0.42256028, -0.33641481},
	// 	{0.2236068, -0.16666667, -0.2191693, -0.35087595, -0.1814372},
	// 	{-0.4472136, 0., 0.08691197, 0.2339173, -0.06047907},
	// 	{0., 0.08333333, -0.13225734, 0.35087595, 0.33641481},
	// 	{0.2236068, 0.375, 0.26451468, 0.39992313, 0.39689388},
	// 	{0.67082039, 0.625, 0.75197744, 0.30560163, 0.45737295},
	// }

}
