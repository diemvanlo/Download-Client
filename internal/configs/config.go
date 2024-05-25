package configs

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type ConfigFilePath string

type Config struct {
	Account  Account  `yaml:"account"`
	Database Database `yaml:"database"`
}

func NewConfig(filePath ConfigFilePath) (Config, error) {
	configBytes, err := os.ReadFile(string(filePath))
	if err != nil {
		return Config{}, fmt.Errorf("Failed to read YAML file: %w", err)
	}

	config := Config{}
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return Config{}, fmt.Errorf("Failed to read YAML file: %w", err)
	}

	return config, nil
}
