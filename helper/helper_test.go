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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelperFunctions(t *testing.T) {
	enabled := true
	update = &enabled
	SaveGoldenData(t, "test", []byte("Testing golden data"))
	data := GetGoldenData(t, "test")
	assert.Equal(t, []byte("Testing golden data"), data,
		"Unexpected error getting the golden data")

	data = GetJSONFile(t, "json_file")
	assert.Equal(t, []byte{}, data, "Expected empty data")

}
