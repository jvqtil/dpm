package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type registryItem struct {
	PkgName     string `json:"pkg_name"`
	Version     string `json:"version"`
	SourceType  string `json:"source_type"`
	Source      string `json:"source"`
	AssetName   string `json:"asset_name"`
	AssetURL    string `json:"asset_url"`
	Binary      string `json:"binary"`
	InstalledAt string `json:"installed_at"`
	LastUpdated string `json:"last_updated"`
}

type registry struct {
	Packages map[string]registryItem `json:"packages"`
}

func getRegPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "state", "dpm", "registry.json")
}

func loadRegistry() (*registry, error) {
	data, err := os.ReadFile(getRegPath())
	if err != nil {
		if os.IsNotExist(err) {
			return &registry{Packages: map[string]registryItem{}}, nil
		}
		return nil, fmt.Errorf("failed to read registry: %w", err)
	}

	var reg registry
	if err := json.Unmarshal(data, &reg); err != nil {
		return nil, fmt.Errorf("failed to parse registry file: %w", err)
	}

	if reg.Packages == nil {
		reg.Packages = map[string]registryItem{}
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
