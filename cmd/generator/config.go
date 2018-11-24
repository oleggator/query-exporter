package main

import (
	"bufio"
	"gopkg.in/yaml.v2"
	"os"
)

// Config struct
type Config struct {
	Host           string `yaml:"host"`
	Port           uint16 `yaml:"port"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	DBName         string `yaml:"db_name"`
	PeopleCount    int    `yaml:"people_count"`
	CountriesCount int    `yaml:"countries_count"`
	CitiesCount    int    `yaml:"cities_count"`
}

func getConfig(configPath string) (*Config, error) {
	configFile, err := os.Open(configPath)
	defer configFile.Close()
	if err != nil {
		return nil, err
	}

	configDecoder := yaml.NewDecoder(bufio.NewReader(configFile))
	config := &Config{}
	err = configDecoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
