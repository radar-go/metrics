// Package config defines the service configurations.
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
	"io/ioutil"

	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
)

// Config stores the service configurations.
type Config struct {
	dir  string
	Port int `yaml:"port"`
}

// New returns a new Config object.
func New(dir string) (*Config, error) {
	conf := &Config{
		dir:  dir,
		Port: 8080,
	}

	err := conf.ReloadConfigs()
	if err != nil {
		return conf, err
	}

	return conf, err
}

// ReloadConfigs reload the configurations stored in the configuration file.
func (c *Config) ReloadConfigs() error {
	yamlFile, err := ioutil.ReadFile(c.dir + "/metrics.yaml")
	if err != nil {
		glog.Errorf("Error reading the config file: %s ", err)
		return err
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		glog.Errorf("Error unmarshaling the configurations: %s ", err)
		return err
	}

	return nil
}
