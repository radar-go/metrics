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
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"
	"github.com/valyala/fasthttp"

	"github.com/radar-go/metrics/pkg/api/controller"
	"github.com/radar-go/metrics/pkg/config"
)

type exiter func(code int)

// API type set up the api server to response to the petitions.
type API struct {
	cfg      *config.Config
	listener net.Listener
	exit     chan os.Signal
	reload   chan os.Signal
	exitfn   exiter
}

// New returns a new API object.
func New() *API {
	flag.Parse()

	cfgDir := os.Getenv("CONF_DIR")
	if cfgDir == "" {
		cfgDir = "./conf"
	}

	glog.Infof("Reading the configurations from %s", cfgDir)
	cfg, err := config.New(cfgDir)
	if err != nil {
		glog.Warningf("Unable to read the configs from %s, trusting the default values",
			cfgDir)
	}

	a := &API{
		cfg:    cfg,
		exit:   make(chan os.Signal, 1),
		reload: make(chan os.Signal, 1),
		exitfn: func(code int) { os.Exit(code) },
	}

	signal.Notify(a.exit, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(a.reload, syscall.SIGUSR1)

	return a
}

func (a *API) Start() error {
	var err error

	err = a.cfg.ReloadConfigs()
	if err != nil {
		glog.Errorf("Error reloading the configurations: %s", err)
		return err
	}

	c := controller.New(a.cfg)
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

// SignalListen starts to listen to the signals to reload and exit the api service.
func (a *API) SignalListen() {
	select {
	case <-a.exit:
		glog.Info("Stoping server...")
		err := a.Stop()
		if err != nil {
			glog.Exitf("Error stoping the server: %s", err)
		}

		a.exitfn(0)
	case <-a.reload:
		glog.Info("Reloading server...")
		err := a.Reload()
		if err != nil {
			glog.Exitf("Error reloading the server: %s", err)
		}

		glog.Info("Server reloaded")
	}
}
