/////////////////////////////////////////////////////////////////////////////////
//
// config.go
//
// Written by Lorenz Breidenbach lob@open.ch, December 2015
// Copyright (c) 2015 Open Systems AG, Switzerland
// All Rights Reserved.
//
/////////////////////////////////////////////////////////////////////////////////

// Package for parsing goprobe config files.
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/els0r/goProbe/pkg/capture"
)

type Config struct {
	DBPath      string                           `json:"db_path"`
	Interfaces  map[string]capture.CaptureConfig `json:"interfaces"`
	SyslogFlows bool                             `json:"syslog_flows"`
	Logging     LogConfig                        `json:"logging"`
}

type LogConfig struct {
	Destination string `json:"destination"`
	Level       string `json:"level"`
}

func New() *Config {
	interfaces := make(map[string]capture.CaptureConfig)

	return &Config{
		Interfaces: interfaces,
		Logging:    LogConfig{"syslog", "info"}, // default config is syslog
	}
}

func (c *Config) Validate() error {
	if c.DBPath == "" {
		return fmt.Errorf("Database path must not be empty")
	}
	for iface, cc := range c.Interfaces {
		err := cc.Validate()
		if err != nil {
			return fmt.Errorf("Interface '%s' has invalid configuration: %s", iface, err)
		}
	}
	return nil
}

func ParseFile(path string) (*Config, error) {
	config := New()

	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}
