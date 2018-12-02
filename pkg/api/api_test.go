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
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

const unexpectedError = "Unexpected error"

func TestAPI(t *testing.T) {
	os.Setenv("CONF_DIR", "./testdata")
	a := New()
	assert.NotNil(t, a, "Error creating the new api object")

	err := a.Start()
	assert.NoError(t, err, unexpectedError)

	err = a.Reload()
	assert.NoError(t, err, unexpectedError)

	err = a.Stop()
	assert.NoError(t, err, unexpectedError)

	err = a.Reload()
	assert.NoError(t, err, unexpectedError)

	err = a.Start()
	assert.EqualError(t, err, "Listener already started", unexpectedError)

	err = a.Stop()
	assert.NoError(t, err, unexpectedError)

	a.reload <- syscall.SIGUSR1
	a.SignalListen()

	a.exitfn = func(code int) {}
	a.exit <- syscall.SIGTERM
	a.SignalListen()
}
