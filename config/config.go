package config

import (
	"fmt"
	"os"

	"github.com/adaggerboy/utasksd/models/config"
	"gopkg.in/yaml.v2"
)

var GlobalConfig config.Config

func Apply(config config.Config, err error) {
	if err != nil {
		panic(err)
	}
	GlobalConfig = config
}

func Load(filename string) (config config.Config, err error) {
	err = nil
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		err = fmt.Errorf("config loading: %s", err)
		return
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		err = fmt.Errorf("config parsing: %s", err)
	}
	return
}
