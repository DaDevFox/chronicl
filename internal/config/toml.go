package config

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

// LoadTOML reads the TOML config file
func LoadTOML(path string) (*Config, error) {
	var config Config
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = toml.Unmarshal(data, &config)
	return &config, err
}

