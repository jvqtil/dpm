package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"time"

	"github.com/dustin/go-humanize"
)

func switchTag(pkg Package, tag Tag, msg string) error {
	reg, err := loadRegistry()
	if err != nil {
		return err
	}

	src, err := resolveAsset(pkg, &tag)
	if err != nil {
		return err
	}

	// Backup existing binary before replacement
	backupPath := pkg.BinaryPath + ".backup"
	if err = cpToDest(pkg.BinaryPath, backupPath); err != nil {
		return err
	}
	defer rmBin(backupPath) // Clean up backup on success

	if err := cpToDest(src, pkg.BinaryPath); err != nil {
		return err
	}

	pkg.AddTag(tag)
	if err := pkg.SetCurrentTag(tag); err != nil {
		// Restore backup on failure
		if _, statErr := os.Stat(backupPath); statErr == nil {
			if err = cpToDest(backupPath, pkg.BinaryPath); err != nil {
				return err
			}
		}
		return err
	}
	pkg.LastUpdated = time.Now().Format(time.RFC3339)
	reg.Packages[pkg.Name] = pkg

	if err := saveRegistry(reg); err != nil {
		// Restore backup on failure
		if _, statErr := os.Stat(backupPath); statErr == nil {
			if err = cpToDest(backupPath, pkg.BinaryPath); err != nil {
				return err
			}
		}
		return err
	}

	fmt.Printf("=> %s\n", msg)
	return nil
}

func fetchTag(name, tagName string) error {
	reg, err := loadRegistry()
	if err != nil {
		return err
	}

	pkg, ok := reg.Packages[name]
	if !ok {
		return fmt.Errorf("package %s is not installed", name)
	}

	var tag *Tag
	switch pkg.SourceType {
	case "github.com":
		release, err := checkGithubTag(pkg.Source, tagName)
		if err != nil {
			return err
		}
		tag, err = resolveGhTag(release)
		if err != nil {
			return err
		}
	default:
		fmt.Println("Fetching is only available for packages from Git hostings")
		return nil
	}

	_, err = resolveAsset(pkg, tag)
	if err != nil {
		return err
	}

	pkg.AddTag(*tag)
	reg.Packages[pkg.Name] = pkg

	if err := saveRegistry(reg); err != nil {
		return err
	}

	fmt.Printf("=> Fetched tag %s for package %s\n", tag.Name, green(pkg.Name))

	return nil
}

func useTag(name, tagName string) error {
	_, pkg, _, err := resolvePkgTag(name, tagName)
	if err != nil {
		return err
	}

	if tagName == "" {
		if len(pkg.Tags) == 0 {
			return fmt.Errorf("no tags found for %s", pkg.Name)
		}
		tags := make([]string, len(pkg.Tags))
		for i, t := range pkg.Tags {
			tags[i] = t.Name
		}
		sort.Strings(tags)
		slices.Reverse(tags)

		tagName, err = picker(tags, fmt.Sprintf("Select tag for %s - current %s", green(pkg.Name), pkg.CurrentTag.Name))
		if err != nil {
			return err
		}
	}

	if pkg.CurrentTag.Name == tagName {
		fmt.Printf("%s is already on %s\n", green(pkg.Name), tagName)
		return nil
	}

	var tag *Tag
	for i, t := range pkg.Tags {
		if t.Name == tagName {
			tag = &pkg.Tags[i]
			break
		}
	}
	if tag == nil {
		return fmt.Errorf("tag %q not found for package %s", tagName, pkg.Name)
	}

	return switchTag(*pkg, *tag, fmt.Sprintf("Switched %s to %s", green(pkg.Name), tagName))
}

func showTagInfo(name, tagName string, jsonOut bool) error {
	_, pkg, tag, err := resolvePkgTag(name, tagName)
	if err != nil {
		return err
	}

	if jsonOut {
		data, err := json.MarshalIndent(tag, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	} else {

		fmt.Println(border)
		fmt.Printf("%s (%s)\n", green(pkg.Name), pkg.SourceType)
		fmt.Println()

		if tag == nil {
			return fmt.Errorf("tag %q not found for package %s", tagName, pkg.Name)
		}
		fmt.Printf("%s\n", tag.Name)
		fmt.Printf("↓ %s %s\n", tag.AssetName, humanize.Bytes(uint64(tag.AssetSize)))
		if pkg.SourceType == "local" {
			fmt.Printf("from:\n  %s\n", tag.AssetPath)
		} else {
			fmt.Printf("from:\n  %s\n", tag.AssetURL)
		}

		fmt.Println(border)
	}
	return nil
}

func removeTag(name, tagName string) error {
	reg, pkg, tag, err := resolvePkgTag(name, tagName)
	if err != nil {
		return err
	}

	if tag == nil {
		return fmt.Errorf("tag %q not found for package %s", tagName, name)
	}

	if pkg.CurrentTag.Name == tag.Name {
		fmt.Printf("Cannot remove current tag (%q), switch to another tag first\n", tag.Name)
		return nil
	}

	if !confirm(fmt.Sprintf("Remove tag %s from package %s?", tag.Name, green(pkg.Name))) {
		return nil
	}

	if err := pkg.RemoveTag(*tag); err != nil {
		return err
	}
	reg.Packages[pkg.Name] = *pkg

	if err := saveRegistry(reg); err != nil {
		return err
	}

	fmt.Printf("=> Removed tag %s from package %s\n", tag.Name, green(pkg.Name))

	if pkg.SourceType != "local" {
		cachePath := filepath.Join(pkg.ResolveCachePath(), tag.Name)
		if _, err := os.Stat(cachePath); err == nil {
			if confirm(fmt.Sprintf("Remove cache for %s?", tag.Name)) {
				if err := os.RemoveAll(cachePath); err != nil {
					return fmt.Errorf("failed to clear cache: %w", err)
				} else {
					fmt.Printf("=> Cleared cache for %s\n", tag.Name)
				}
			}
		}
	}

	return nil
}
