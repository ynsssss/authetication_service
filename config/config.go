package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Url string `yaml:"url"`
	} `yaml:"database"`
	Secret string `yaml:"secret"`
}

func Get(filePath string) (*Config, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var config *Config
	err = yaml.NewDecoder(f).Decode(&config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func MustGet(filePath string) *Config {
	cfg, err := Get(filePath)
	if err != nil {
		panic(err)
	}
	return cfg
}
