package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	dir, _ := os.Getwd()
	paths := []string{
		os.Getenv("HOME"),
		dir,
	}

	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			fmt.Printf("%v", err)
		}
		defer f.Close()

		files, err := f.Readdir(-1)
		if err != nil {
			fmt.Printf("%v", err)
		}

		for _, file := range files {
			if strings.HasPrefix(file.Name(), "chronicl") {
				switch filepath.Ext(file.Name()) {
				case ".yaml", ".yml", "":
					return LoadYAML(filepath.Join(path, file.Name()))
				case ".toml":
					return LoadTOML(filepath.Join(path, file.Name()))
					// case ".textproto":
					// 	return LoadTextProto(path)
				}
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
