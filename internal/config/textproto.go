package config

import (
	"os"

	"google.golang.org/protobuf/encoding/prototext"
)

// LoadTextProto reads the textproto config file
func LoadTextProto(path string) (*Config, error) {
	var config Config
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = prototext.Unmarshal(data, &config)
	return &config, err
}

