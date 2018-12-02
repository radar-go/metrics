package api

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

func TestAPI(t *testing.T) {
	a := New("./testdata")
	assert.NotNil(t, a, "Error creating the new api object")

	err := a.Start()
	assert.NoError(t, err, "Unexpected error")

	err = a.Reload()
	assert.NoError(t, err, "Unexpected error")

	err = a.Stop()
	assert.NoError(t, err, "Unexpected error")

	err = a.Reload()
	assert.NoError(t, err, "Unexpected error")

	err = a.Start()
	assert.EqualError(t, err, "Listener already started", "Unexpected error")

	err = a.Stop()
	assert.NoError(t, err, "Unexpected error")
}
