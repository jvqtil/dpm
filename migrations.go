package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// This is a migration suite, which is going to be removed in v1.0.0
func runMigration(migration string) error {
	switch migration {
	case "v0.0.5":
		fmt.Println("In v0.0.5 there was a paths change")
		if !confirm("Do you want to run a migration?") {
			return nil
		}

		err := migrateAfterPathsChange()
		if err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}

		fmt.Println("Successfully migrated!")
	default:
		fmt.Printf("Migration %s not found\n", migration)
	}
	return nil
}

// Paths changed in v0.0.5, this is a migration for it
func migrateAfterPathsChange() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Join(configDir, "dpm"), 0755); err != nil {
		return err
	}

	prevConfigPath := filepath.Join(homeDir, ".config", "dpm", "config.toml")
	prevRegPath := filepath.Join(homeDir, ".local", "state", "dpm", "registry.json")
	prevCacheDir := filepath.Join(homeDir, ".cache", "dpm")

	currentConfigPath := filepath.Join(configDir, "dpm", "config.toml")
	currentRegPath := filepath.Join(configDir, "dpm", "registry.json")
	currentCacheDir := filepath.Join(cacheDir, "dpm")

	// config migration
	if _, err := os.Stat(currentConfigPath); os.IsNotExist(err) {
		if _, err := os.Stat(prevConfigPath); err == nil {
			if err := os.Rename(prevConfigPath, currentConfigPath); err != nil {
				return fmt.Errorf("failed to migrate config: %w", err)
			}
		}
	}

	// registry migration
	if _, err := os.Stat(currentRegPath); os.IsNotExist(err) {
		if _, err := os.Stat(prevRegPath); err == nil {
			if err := os.Rename(prevRegPath, currentRegPath); err != nil {
				return fmt.Errorf("failed to migrate registry:%w", err)
			}
		}
	}

	// cache migration
	if _, err := os.Stat(currentCacheDir); os.IsNotExist(err) {
		if _, err := os.Stat(prevCacheDir); err == nil {
			if err := os.Rename(prevCacheDir, currentCacheDir); err != nil {
				return fmt.Errorf("failed to migrate cache:%w", err)
			}
		}
	}
	return nil
}
