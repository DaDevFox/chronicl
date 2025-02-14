package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config structure
type Config struct {
	CommitTypes []string `json:"commit_types" toml:"commit_types" yaml:"commit_types"`
	Scopes      []string `json:"scopes" toml:"scopes" yaml:"scopes"`
	AutoCommit  bool     `json:"auto_commit" toml:"auto_commit" yaml:"auto_commit"`
}

// LoadConfig tries to read config from multiple sources
func LoadConfig() (*Config, error) {
	var config Config
	paths := []string{
		filepath.Join(os.Getenv("HOME"), ".rgrc"), // YAML default
		filepath.Join(os.Getenv("HOME"), ".rgrc.toml"),
		filepath.Join(os.Getenv("HOME"), ".rgrc.textproto"),
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			switch filepath.Ext(path) {
			case ".yaml", ".yml", "":
				return LoadYAML(path)
			case ".toml":
				return LoadTOML(path)
			case ".textproto":
				return LoadTextProto(path)
			}
		}
	}

	// Default config if no file found
	return &Config{
		CommitTypes: []string{"feat", "fix", "docs", "style", "refactor", "perf", "test", "chore", "ci"},
		Scopes:      []string{},
		AutoCommit:  false,
	}, nil
}

