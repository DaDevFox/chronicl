package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// TODO: allow custom format strings for actual message, not just display helper
type MessageFormatterConfig struct {
	Key              string `json:"key" toml:"key" yaml:"key"`
	Description      string `json:"description" toml:"description" yaml:"description"`
	HelperFormatting string `json:"hint_format" toml:"hint_format" yaml:"hint_format"` // formatting: {scope} for commit scope, {message} for commit message text
	// TODO: add {author} substitution
}

// Config structure
type Config struct {
	CommitTypes             []MessageFormatterConfig            `json:"commit_types" toml:"commit_types" yaml:"commit_types"`
	Scopes                  []MessageFormatterConfig            `json:"scopes" toml:"scopes" yaml:"scopes"`
	CommitTypeToAnnotations map[string][]MessageFormatterConfig `json:"commit_annotations" toml:"commit_annotations" yaml:"commit_annotations"`
	ScopeToAnnotations      map[string][]MessageFormatterConfig `json:"scope_annotations" toml:"scope_annotations" yaml:"scope_annotations"`
	AutoCommit              bool                                `json:"auto_commit" toml:"auto_commit" yaml:"auto_commit"`
	PadCommitType           bool                                `json:"pad_commit_type" toml:"pad_commit_type" yaml:"pad_commit_type"`
	PadScope                bool                                `json:"pad_scope" toml:"pad_scope" yaml:"pad_scope"`
	AbortOnEmptyCommit      bool                                `json:"abort_on_empty_commit" toml:"abort_on_empty_commit" yaml:"abort_on_empty_commit"`
	PadCharacter            string                              `json:"pad_character" toml:"pad_character" yaml:"pad_character"`
	CommitCommandFormat     string                              `json:"commit_command_format" toml:"commit_command_format" yaml:"commit_command_format"`
}

// LoadConfig tries to read config from multiple sources
func LoadConfig() (*Config, error) {
	// var config Config
	dir, _ := os.Getwd()
	paths := []string{
		dir,
		os.Getenv("HOME"),
		path.Join(os.Getenv("HOME"), ".config/chronicl"),
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
				case ".yaml", ".yml":
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
		CommitTypes:         []MessageFormatterConfig{{"feat", "semver MINOR", ""}, {"fix", "semver PATH", ""}, {"docs", "docs", ""}, {"style", "style", ""}, {"refactor", "refactor", ""}, {"perf", "perf", ""}, {"test", "test", ""}, {"chore", "chore", ""}, {"ci", "ci", ""}},
		Scopes:              []MessageFormatterConfig{},
		AutoCommit:          false,
		CommitCommandFormat: "",
	}, nil
}
