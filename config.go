package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type config struct {
	BinDir           string `toml:"bin_dir"`
	CacheDir         string `toml:"cache_dir"`
	AssumeSourceType string `toml:"assume_source_type"`
}

var cfg config

func defaultConfig() config {
	homeDir, _ := os.UserHomeDir()
	return config{
		BinDir:           "/usr/local/bin",
		CacheDir:         filepath.Join(homeDir, ".cache", "dpm"),
		AssumeSourceType: "github.com",
	}
}

func initConfig() error {
	homeDir, _ := os.UserHomeDir()
	path := filepath.Join(homeDir, ".config", "dpm", "config.toml")

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			cfg = defaultConfig()
			return nil
		}
		return fmt.Errorf("failed to read config: %w", err)
	}

	cfg = defaultConfig()
	if _, err := toml.Decode(string(data), &cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	return nil
}
