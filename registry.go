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

// Legacy schema for migration
type legacyRegistryItem struct {
	Name        string `json:"pkg_name"`
	Version     string `json:"version"`
	AssetName   string `json:"asset_name"`
	AssetURL    string `json:"asset_url"`
	AssetPath   string `json:"asset_path"`
	AssetSize   int64  `json:"asset_size"`
	SourceType  string `json:"source_type"`
	Source      string `json:"source"`
	BinaryPath  string `json:"binary_path"`
	InstalledAt string `json:"installed_at"`
	LastUpdated string `json:"last_updated"`
}

type legacyRegistry struct {
	Packages map[string]legacyRegistryItem `json:"packages"`
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
		// Try migrating from legacy schema
		var legacyReg legacyRegistry
		if legacyErr := json.Unmarshal(data, &legacyReg); legacyErr == nil && legacyReg.Packages != nil {
			reg.Packages = make(map[string]pkg)
			for name, item := range legacyReg.Packages {
				// Migrate legacy item to current schema
				currentTag := tag{
					TagName:   item.Version,
					AssetName: item.AssetName,
					AssetURL:  item.AssetURL,
					AssetPath: item.AssetPath,
					AssetSize: item.AssetSize,
				}
				reg.Packages[name] = pkg{
					Name:        item.Name,
					CurrentTag:  currentTag,
					Tags:        []tag{currentTag},
					SourceType:  item.SourceType,
					Source:      item.Source,
					BinaryPath:  item.BinaryPath,
					InstalledAt: item.InstalledAt,
					LastUpdated: item.LastUpdated,
				}
			}
		} else {
			return nil, fmt.Errorf("failed to parse registry file: %w", err)
		}
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
