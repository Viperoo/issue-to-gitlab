package main

import (
	"os"

	"github.com/BurntSushi/toml"
)

type ConfigStructure struct {
	Host              string
	Token             string
	DefaultAssigneeId int
}

var Config = ConfigStructure{}

func ReadConfig(configfile string) {
	_, err := os.Stat(configfile)
	if err != nil {
		logger.Criticalf("Configuration file: %s not found", configfile)
		os.Exit(1)
	}

	var config ConfigStructure
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		logger.Critical(err)
	}

	Config = config
}
