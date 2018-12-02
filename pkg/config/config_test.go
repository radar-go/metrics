package config

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

func TestConfig(t *testing.T) {
	cfg, err := New("../../conf")
	assert.NoError(t, err, "Unexpected error reading the configurations")
	assert.NotNil(t, cfg, "Expected the configs to be created")
	assert.Equal(t, 6000, cfg.Port, "Unexpected port in configs")
	assert.Equal(t, "../../conf", cfg.dir, "Unexpected directory in configs")

	cfg, err = New("./testdata")
	assert.Error(t, err, "Expected error reading the configs")
	assert.EqualError(t, err,
		"yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `6000` into int",
		"Unexpected error")
	assert.Equal(t, 8080, cfg.Port, "Unexpected port in configs")
	assert.Equal(t, "./testdata", cfg.dir, "Unexpected directory in configs")

	cfg, err = New("./wrong_dir")
	assert.Error(t, err, "Expected error reading the configs")
	assert.EqualError(t, err,
		"open ./wrong_dir/metrics.yaml: no such file or directory",
		"Unexpected error")
	assert.Equal(t, 8080, cfg.Port, "Unexpected port in configs")
	assert.Equal(t, "./wrong_dir", cfg.dir, "Unexpected directory in configs")
}
