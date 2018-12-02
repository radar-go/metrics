// Package api implements the service api basic functionality, that includes
// creating and initializing the api and starting, stoping and reloading it.
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
	"fmt"
	"net"
	"time"

	"github.com/golang/glog"
	"github.com/valyala/fasthttp"

	"github.com/radar-go/metrics/pkg/api/controller"
	"github.com/radar-go/metrics/pkg/config"
)

// API type set up the api server to response to the petitions.
type API struct {
	cfg      *config.Config
	listener net.Listener
}

// New returns a new API object.
func New(cfgDir string) *API {
	glog.Infof("Reading the configurations from %s", cfgDir)
	cfg, err := config.New(cfgDir)
	if err != nil {
		glog.Warningf("Unable to read the configs from %s, trusting the default values",
			cfgDir)
	}

	return &API{
		cfg: cfg,
	}
}

// Start the API server.
func (a *API) Start() error {
	var err error

	a.cfg.ReloadConfigs()
	c := controller.New(a.cfg)
	if err != nil {
		glog.Errorf("Error creating the controller: %s", err)
		return err
	}

	server := fasthttp.Server{
		Handler:           fasthttp.CompressHandler(c.Router.Handler),
		ReadBufferSize:    1024 * 64,
		WriteBufferSize:   1024 * 64,
		ReduceMemoryUsage: true,
	}

	if a.listener != nil {
		err = fmt.Errorf("Listener already started")
		glog.Error(err)
		return err
	}

	a.listener, err = net.Listen("tcp4", fmt.Sprint(":", a.cfg.Port))
	if err != nil {
		glog.Errorf("Error creating the listener: %s", err)
		return err
	}

	glog.Infof("Starting server on port %d", a.cfg.Port)
	go func() {
		err := server.Serve(a.listener)
		if err != nil {
			glog.Errorf("Error starting the server: %s", err)
		}
	}()

	return err
}

// Stop the API server.
func (a *API) Stop() error {
	var err error

	if a.listener != nil {
		glog.Info("Stopping the API")
		err = a.listener.Close()
		time.Sleep(time.Second)
		a.listener = nil
	}

	return err
}

// Reload the API server.
func (a *API) Reload() error {
	err := a.Stop()
	if err != nil {
		glog.Errorf("Error stoping the server: %s", err)
		return err
	}

	err = a.Start()
	if err != nil {
		glog.Errorf("Error starting the server: %s", err)
	}

	return err
}
