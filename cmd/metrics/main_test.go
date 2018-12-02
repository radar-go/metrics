package main

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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestService(t *testing.T) {
	os.Setenv("CONF_DIR", "./testdata")
	go main()

	time.Sleep(2 * time.Second)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://localhost:6000/healthcheck")
	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	client.Do(req, resp)
	bodyBytes := resp.Body()
	assert.Equal(t, []byte(`{"status": "ok"}`), bodyBytes)

	reload <- syscall.SIGUSR1
	time.Sleep(2 * time.Second)
	client.Do(req, resp)
	bodyBytes = resp.Body()
	assert.Equal(t, []byte(`{"status": "ok"}`), bodyBytes)

	exit <- syscall.SIGTERM
	time.Sleep(2 * time.Second)
	client.Do(req, resp)
	bodyBytes = resp.Body()
	assert.Equal(t, []byte(""), bodyBytes)
}
