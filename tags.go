package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

func switchTag(pkg Package, tag Tag, msg string) error {
	src, err := resolveAsset(pkg, &tag)
	if err != nil {
		return err
	}

	if err := cpToDest(src, pkg.BinaryPath); err != nil {
		return err
	}

	reg, err := loadRegistry()
	if err != nil {
		return err
	}

	pkg.AddTag(tag)
	if err := pkg.SetCurrentTag(tag); err != nil {
		return err
	}
	pkg.LastUpdated = time.Now().Format(time.RFC3339)
	reg.Packages[pkg.Name] = pkg

	if err := saveRegistry(reg); err != nil {
		return err
	}

	fmt.Printf("=> %s\n", msg)
	return nil
}

func fetchTag(name, tag string) error {
	reg, err := loadRegistry()
	if err != nil {
		return err
	}

	pkg, ok := reg.Packages[strings.ToLower(name)]
	if !ok {
		return fmt.Errorf("package %s is not installed", name)
	}

	var targetTag *Tag
	switch pkg.SourceType {
	case "github.com":
		release, err := checkGithubTag(pkg.Source, tag)
		if err != nil {
			return err
		}
		targetTag, err = resolveGhTag(release)
		if err != nil {
			return err
		}
	default:
		fmt.Println("Fetching is only available for packages from Git hostings")
		return nil
	}

	_, err = resolveAsset(pkg, targetTag)
	if err != nil {
		return err
	}

	pkg.AddTag(*targetTag)
	reg.Packages[pkg.Name] = pkg

	if err := saveRegistry(reg); err != nil {
		return err
	}

	fmt.Printf("=> Fetched tag %s for package %s\n", targetTag.Name, green(pkg.Name))

	return nil
}

func useTag(name, tag string) error {
	reg, err := loadRegistry()
	if err != nil {
		return err
	}

	pkg, ok := reg.Packages[strings.ToLower(name)]
	if !ok {
		return fmt.Errorf("package %s is not installed", name)
	}

	if tag == "" {
		if len(pkg.Tags) == 0 {
			return fmt.Errorf("no tags found for %s", pkg.Name)
		}
		tags := make([]string, len(pkg.Tags))
		for i, t := range pkg.Tags {
			tags[i] = t.Name
		}
		sort.Strings(tags)
		slices.Reverse(tags)

		tag, err = picker(tags, fmt.Sprintf("Select tag for %s - current %s", green(pkg.Name), pkg.CurrentTag.Name))
		if err != nil {
			return err
		}
	}

	if pkg.CurrentTag.Name == tag {
		fmt.Printf("%s is already on %s\n", green(pkg.Name), tag)
		return nil
	}

	var target *Tag
	for i, t := range pkg.Tags {
		if t.Name == tag {
			target = &pkg.Tags[i]
			break
		}
	}
	if target == nil {
		return fmt.Errorf("tag %q not found for package %s", tag, pkg.Name)
	}

	return switchTag(pkg, *target, fmt.Sprintf("Switched %s to %s", green(pkg.Name), tag))
}

func showTagInfo(name, tag string, jsonOut bool) error {
	reg, err := loadRegistry()
	if err != nil {
		return err
	}

	pkg, ok := reg.Packages[strings.ToLower(name)]
	if !ok {
		return fmt.Errorf("package %s is not installed", name)
	}

	var target *Tag
	for i, t := range pkg.Tags {
		if t.Name == tag {
			target = &pkg.Tags[i]
			break
		}
	}

	if jsonOut {
		data, err := json.MarshalIndent(target, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	} else {

		fmt.Println(border)
		fmt.Printf("%s (%s)\n", green(pkg.Name), pkg.SourceType)
		fmt.Println()

		if target == nil {
			return fmt.Errorf("tag %q not found for package %s", tag, pkg.Name)
		}
		fmt.Printf("%s\n", target.Name)
		fmt.Printf("↓ %s %s\n", target.AssetName, humanize.Bytes(uint64(target.AssetSize)))
		if pkg.SourceType == "local" {
			fmt.Printf("from:\n  %s\n", target.AssetPath)
		} else {
			fmt.Printf("from:\n  %s\n", target.AssetURL)
		}

		fmt.Println(border)
	}
	return nil
}

func removeTag(name, tag string) error {
	reg, err := loadRegistry()
	if err != nil {
		return err
	}

	pkg, ok := reg.Packages[strings.ToLower(name)]
	if !ok {
		return fmt.Errorf("package %s is not installed", name)
	}

	var target *Tag
	for i, t := range pkg.Tags {
		if t.Name == tag {
			target = &pkg.Tags[i]
			break
		}
	}

	if target == nil {
		return fmt.Errorf("tag %q not found for package %s", tag, name)
	}

	if pkg.CurrentTag.Name == target.Name {
		fmt.Printf("Cannot remove current tag (%q), switch to another tag first\n", tag)
		return nil
	}

	if !confirm(fmt.Sprintf("Remove tag %s from package %s?", target.Name, green(pkg.Name))) {
		return nil
	}

	if err := pkg.RemoveTag(*target); err != nil {
		return err
	}
	reg.Packages[pkg.Name] = pkg

	if err := saveRegistry(reg); err != nil {
		return err
	}

	fmt.Printf("=> Removed tag %s from package %s\n", tag, green(pkg.Name))

	cachePath := filepath.Join(pkg.ResolveCachePath(), target.Name)
	if _, err := os.Stat(cachePath); err == nil {
		if confirm(fmt.Sprintf("Remove cache for %s?", tag)) {
			if err := os.RemoveAll(cachePath); err != nil {
				return fmt.Errorf("failed to clear cache: %w", err)
			} else {
				fmt.Printf("=> Cleared cache for %s\n", tag)
			}
		}
	}

	return nil
}
