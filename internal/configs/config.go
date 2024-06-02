package configs

import (
	"fmt"
	"goload/configs"
	"gopkg.in/yaml.v3"
	"os"
)

type ConfigFilePath string

type Config struct {
	GRPC     GRPC     `yaml:"GRPC"`
	HTTP     HTTP     `yaml:"HTTP"`
	Log      Log      `yaml:"log"`
	Auth     Auth     `yaml:"auth"`
	Database Database `yaml:"database"`
	Cache    Cache    `yaml:"cache"`
}

func NewConfig(filePath ConfigFilePath) (Config, error) {
	var (
		configBytes = configs.DefaultConfigBytes
		config      = Config{}
		err         error
	)

	if filePath != "" {
		configBytes, err = os.ReadFile(string(filePath))
		if err != nil {
			return Config{}, fmt.Errorf("Failed to read YAML file: %w", err)
		}
	}

	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return Config{}, fmt.Errorf("Failed to read YAML file: %w", err)
	}

	return config, nil
}
