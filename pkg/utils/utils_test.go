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
