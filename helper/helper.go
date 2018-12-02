// Package helper implements helper functions for the tests.
package helper

/* Copyright (C) 2018 Radar team (see AUTHORS)

   This file is part of metrics.

   metrics is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   metrics is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with metrics. If not, see <http://www.gnu.org/licenses/>.
*/

import (
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "update .golden files")

const testDir = "testdata"

// SaveGoldenData saves test data in a golden file.
func SaveGoldenData(t *testing.T, name string, data []byte) {
	t.Helper()
	golden := filepath.Join(testDir, name+".golden")
	if *update {
		err := ioutil.WriteFile(golden, data, 0644)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
	}
}

// GetGoldenData gets data from a golden file.
func GetGoldenData(t *testing.T, name string) []byte {
	t.Helper()
	golden := filepath.Join(testDir, name+".golden")
	expected, err := ioutil.ReadFile(golden)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	return expected
}

// GetJSONFile obtain the content of a json file from the testdata directory.
func GetJSONFile(t *testing.T, name string) []byte {
	t.Helper()
	filename := filepath.Join(testDir, name+".json")
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("Unexpected error reading the json file: %s", err)
	}

	return content
}
