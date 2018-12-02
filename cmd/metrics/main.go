// Package main implements the metrics command line.
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
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"

	"github.com/radar-go/metrics/pkg/api"
)

var exit = make(chan os.Signal, 1)
var reload = make(chan os.Signal, 1)

func main() {
	/* Parse the arguments. */
	flag.Parse()

	cfgDir := os.Getenv("CONF_DIR")
	if cfgDir == "" {
		cfgDir = "./conf"
	}

	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(reload, syscall.SIGUSR1)

	a := api.New(cfgDir)

	err := a.Start()
	if err != nil {
		glog.Exitf("Error starting the server: %s", err)
	}

	for {
		select {
		case <-exit:
			glog.Info("Stoping server...")
			err := a.Stop()
			if err != nil {
				glog.Exitf("Error stoping the server: %s", err)
			}

			os.Exit(0)
		case <-reload:
			glog.Info("Reloading server...")
			err := a.Reload()
			if err != nil {
				glog.Exitf("Error reloading the server: %s", err)
			}

			glog.Info("Server reloaded")
		}
	}
}
