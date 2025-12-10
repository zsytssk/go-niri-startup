package config

import (
	"fmt"
	"niri-startup/utils"
	"os"
	"sync"

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

var (
	config Config
	once   sync.Once
)

func GetConfig() (*Config, error) {

	var err error
	once.Do(func() {
		var configPath string
		configPath, err = utils.GetCurDirFileName("config.yml")
		if err != nil {
			return
		}

		configPath, err = utils.GetCurDirFilePath(configPath)
		if err != nil {
			return
		}
		config = Config{}

		configFile, err := os.ReadFile(configPath)
		if err != nil {
			configFile, err = os.ReadFile("config.local.yml")
		}
		if err != nil {
			return
		}
		err = yaml.Unmarshal(configFile, &config)
		if err != nil {
			return
		}
	})
	if err != nil {
		return nil, err
	}
	return &config, err
}

func GetSpadConfig(name string) (*Spad, error) {
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}
	cur, ok := config.SpadMap[name]
	if !ok {
		return nil, fmt.Errorf("can't find spad config name %s", name)
	}
	return &cur, nil
}
