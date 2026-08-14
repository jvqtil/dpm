package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type pkg struct {
	Name        string `json:"pkg_name"`
	CurrentTag  tag    `json:"current_tag"`
	Tags        []tag  `json:"tags"`
	SourceType  string `json:"source_type"`
	Source      string `json:"source"`
	BinaryPath  string `json:"binary_path"`
	InstalledAt string `json:"installed_at"`
	LastUpdated string `json:"last_updated"`
}

type tag struct {
	TagName   string `json:"tag_name"`
	AssetName string `json:"asset_name"`
	AssetURL  string `json:"asset_url"`
	AssetPath string `json:"asset_path"`
	AssetSize int64  `json:"asset_size"`
}

type registry struct {
	Packages map[string]pkg `json:"packages"`
}

func getRegPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "state", "dpm", "registry.json")
}

func loadRegistry() (*registry, error) {
	data, err := os.ReadFile(getRegPath())
	if err != nil {
		if os.IsNotExist(err) {
			return &registry{Packages: map[string]pkg{}}, nil
		}
		return nil, fmt.Errorf("failed to read registry: %w", err)
	}

	var reg registry
	if err := json.Unmarshal(data, &reg); err != nil {
		return nil, fmt.Errorf("failed to parse registry file: %w", err)
	}

	if reg.Packages == nil {
		reg.Packages = map[string]pkg{}
	}

	return &reg, nil
}

func saveRegistry(reg *registry) error {
	path := getRegPath()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create registry dir: %w", err)
	}

	data, err := json.MarshalIndent(reg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal registry: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write registry: %w", err)
	}

	return nil
}
