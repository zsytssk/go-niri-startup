package config

import (
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Spad struct {
	Cmd    string `yaml:"cmd"`
	AppId  string `yaml:"appId"`
	Height int    `yaml:"height"`
	Width  int    `yaml:"width"`
}

type Config struct {
	SpadMap map[string]Spad `yaml:"spadMap"`
}

var config *Config

func GetConfig() (Config, error) {
	if config != nil {
		return *config, nil
	}
	config = &Config{}
	ex, err := os.Executable()
	if err != nil {
		return *config, err
	}
	exDir := filepath.Dir(ex)
	configPath := path.Join(exDir, "config.yaml")

	configFile, err := os.ReadFile(configPath)
	if err != nil {
		configFile, err = os.ReadFile("config.local.yml")
	}
	if err != nil {
		return *config, err
	}
	err = yaml.Unmarshal(configFile, config)
	if err != nil {
		return *config, err
	}
	return *config, err
}
