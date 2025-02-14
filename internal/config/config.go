package config

import (
	"os"
	"path/filepath"
)

type CommitType struct {
	Key         string `json:"key" toml:"key" yaml:"key"`
	Description string `json:"description" toml:"description" yaml:"description"`
}

// Config structure
type Config struct {
	CommitTypes   []CommitType `json:"commit_types" toml:"commit_types" yaml:"commit_types"`
	Scopes        []string     `json:"scopes" toml:"scopes" yaml:"scopes"`
	AutoCommit    bool         `json:"auto_commit" toml:"auto_commit" yaml:"auto_commit"`
	PadCommitType bool         `json:"pad_commit_type" toml:"pad_commit_type yaml:"pad_commit_type"`
	PadScope      bool         `json:"pad_scope" toml:"pad_scope" yaml:"pad_scope"`
}

// LoadConfig tries to read config from multiple sources
func LoadConfig() (*Config, error) {
	// var config Config
	paths := []string{
		filepath.Join(os.Getenv("HOME"), ".chroniclrc"), // YAML default
		filepath.Join(os.Getenv("HOME"), ".chroniclrc.toml"),
		// filepath.Join(os.Getenv("HOME"), ".rgrc.textproto"),
	}

	for _, path := range paths {
		if stat, err := os.Stat(path); err == nil && stat.Name() == "chronicl" {
			switch filepath.Ext(path) {
			case ".yaml", ".yml", "":
				return LoadYAML(path)
			case ".toml":
				return LoadTOML(path)
				// case ".textproto":
				// 	return LoadTextProto(path)
			}
		}
	}

	// Default config if no file found
	return &Config{
		CommitTypes: []CommitType{{"feat", "semver MINOR"}, {"fix", "semver PATH"}, {"docs", "docs"}, {"style", "style"}, {"refactor", "refactor"}, {"perf", "perf"}, {"test", "test"}, {"chore", "chore"}, {"ci", "ci"}},
		Scopes:      []string{},
		AutoCommit:  false,
	}, nil
}
