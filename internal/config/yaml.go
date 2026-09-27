package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadYAML reads the YAML config file
func LoadYAML(path string) (*Config, error) {
	var config Config
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("tried to read %s: %w", path, err)
	}
	err = yaml.Unmarshal(data, &config)
	return &config, err
}
