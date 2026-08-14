package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func updatePkg(existingPkg pkg, newTag tag) error {
	var src string
	tmp := &pkg{
		Name:       existingPkg.Name,
		SourceType: existingPkg.SourceType,
		Source:     existingPkg.Source,
		CurrentTag: newTag,
	}
	var err error
	src, err = resolveAsset(tmp)
	if err != nil {
		return err
	}
	if tmp.CurrentTag.AssetPath != "" {
		newTag.AssetPath = tmp.CurrentTag.AssetPath
		newTag.AssetSize = tmp.CurrentTag.AssetSize
	} else {
		newTag.AssetPath = src
	}

	dest := existingPkg.BinaryPath

	backupPath := filepath.Join(os.TempDir(), "dpm-backup-"+filepath.Base(dest))
	if _, err := os.Stat(dest); err == nil {
		data, err := os.ReadFile(dest)
		if err != nil {
			return fmt.Errorf("failed to read existing binary for backup: %w", err)
		}
		if err := os.WriteFile(backupPath, data, 0755); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
		defer os.Remove(backupPath)
	}

	if err := cpToDest(src, dest); err != nil {
		if _, statErr := os.Stat(backupPath); statErr == nil {
			if restoreErr := cpToDest(backupPath, dest); restoreErr != nil {
				return fmt.Errorf("copy failed and restoration failed: copy error: %w, restore error: %v", err, restoreErr)
			}
			return fmt.Errorf("copy failed, rolled back: %w", err)
		}
		return err
	}

	reg, err := loadRegistry()
	if err != nil {
		if _, statErr := os.Stat(backupPath); statErr == nil {
			if restoreErr := cpToDest(backupPath, dest); restoreErr != nil {
				return fmt.Errorf("failed to load registry and restoration failed: load error: %w, restore error: %v", err, restoreErr)
			}
			return fmt.Errorf("failed to load registry, rolled back: %w", err)
		}
		rmBin(dest)
		return fmt.Errorf("failed to load registry: %w", err)
	}

	updatedPkg := existingPkg
	updatedPkg.CurrentTag = newTag
	updatedPkg.LastUpdated = time.Now().Format(time.RFC3339)

	found := -1
	for i, t := range updatedPkg.Tags {
		if t.TagName == newTag.TagName {
			found = i
			break
		}
	}
	if found != -1 {
		updatedPkg.Tags[found] = newTag
	} else {
		updatedPkg.Tags = append(updatedPkg.Tags, newTag)
	}

	reg.Packages[updatedPkg.Name] = updatedPkg

	if err := saveRegistry(reg); err != nil {
		if _, statErr := os.Stat(backupPath); statErr == nil {
			if restoreErr := os.Rename(backupPath, dest); restoreErr != nil {
				return fmt.Errorf("failed to save registry and restoration failed: save error: %w, restore error: %v", err, restoreErr)
			}
			return fmt.Errorf("failed to save registry, rolled back: %w", err)
		}
		return fmt.Errorf("failed to save registry: %w", err)
	}

	fmt.Printf("=> Updated package %s from %s to %s\n", green(existingPkg.Name), existingPkg.CurrentTag.TagName, newTag.TagName)
	return nil
}

func updateTarget(pkgName string) error {
	reg, err := loadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	existingPkg, ok := reg.Packages[strings.ToLower(pkgName)]
	if !ok {
		return fmt.Errorf("package %s is not installed", pkgName)
	}

	var newPkg *pkg
	switch existingPkg.SourceType {
	case "local":
		fmt.Printf("=> %s is a local package. Local packages can't be updated. Please reinstall it manually\n", green(pkgName))
		return nil
	case "direct":
		fmt.Printf("=> %s was installed from direct URL. dpm can't fetch updates for it. Please reinstall it manually\n", green(pkgName))
		return nil
	case "github.com":
		release, err := checkGithubTag(existingPkg.Source, "")
		if err != nil {
			return fmt.Errorf("failed to check version: %w", err)
		}

		if existingPkg.CurrentTag.TagName == release.TagName {
			fmt.Printf("%s is already up to date (%s)\n", green(pkgName), existingPkg.CurrentTag.TagName)
			return nil
		}

		newPkg, err = resolveGithub(existingPkg.Source, existingPkg.Name, release)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported source type %s", existingPkg.SourceType)
	}

	return updatePkg(existingPkg, newPkg.CurrentTag)
}

func updateAll() error {
	reg, err := loadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	if len(reg.Packages) == 0 {
		fmt.Println("No packages installed")
		return nil
	}

	type pending struct {
		item    pkg
		release *ghRelease
	}

	var updates []pending

	fmt.Println("Checking for updates, might take some time...")
	for _, p := range reg.Packages {
		if p.SourceType == "github.com" {
			release, err := fetchGhRelease(p.Source, "")
			if err != nil {
				fmt.Printf("Skipping %s: %v\n", p.Name, err)
				continue
			}

			if release.TagName != p.CurrentTag.TagName {
				updates = append(updates, pending{
					item:    p,
					release: release,
				})
			}
		}
	}

	if len(updates) == 0 {
		fmt.Println("Everything is up to date!")
		return nil
	}

	fmt.Printf("\n=> Updates available (%d)\n", len(updates))
	for _, u := range updates {
		fmt.Printf("  %s  %s -> %s\n", u.item.Name, u.item.CurrentTag.TagName, u.release.TagName)
	}

	if !confirm("\nUpdate all?") {
		return nil
	}

	for _, u := range updates {
		fmt.Println()
		newPkg, err := resolveGithub(u.item.Source, u.item.Name, u.release)
		if err != nil {
			return err
		}
		if err := updatePkg(u.item, newPkg.CurrentTag); err != nil {
			fmt.Printf("failed to update %s: %v\n", u.item.Name, err)
		}
	}

	return nil
}
