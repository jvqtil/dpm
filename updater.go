package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

func update(i registryItem, pkg *pkg) error {
	fmt.Printf("=> Updating package %s from %s to %s\n", green(pkg.Name), i.Version, pkg.Version)

	src := pkg.Source
	if pkg.SourceType != "local" {
		var err error
		src, err = downloadAsset(pkg)
		if err != nil {
			return err
		}
	}

	dest := filepath.Join(cfg.BinDir, pkg.Name)
	if err := cpToDest(src, dest); err != nil {
		return err
	}

	reg, err := loadRegistry()
	if err != nil {
		rmBin(dest)
		return fmt.Errorf("failed to load registry: %w", err)
	}

	reg.Packages[pkg.Name] = registryItem{
		PkgName:     pkg.Name,
		Version:     pkg.Version,
		SourceType:  pkg.SourceType,
		Source:      pkg.Source,
		AssetName:   pkg.AssetName,
		AssetURL:    pkg.AssetURL,
		Binary:      dest,
		InstalledAt: i.InstalledAt,
		LastUpdated: time.Now().Format("02 Jan 06 15:04"),
	}

	if err := saveRegistry(reg); err != nil {
		return fmt.Errorf("failed to save registry, rolled back: %w", err)
	}

	fmt.Printf("=> Updated package %s from %s to %s\n", green(pkg.Name), i.Version, pkg.Version)
	return nil
}

func updateTarget(pkgName string) error {
	reg, err := loadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	i, ok := reg.Packages[strings.ToLower(pkgName)]
	if !ok {
		return fmt.Errorf("package %s is not installed", pkgName)
	}

	if i.SourceType == "local" {
		fmt.Printf("=> %s is a local package. Local packages can't be updated. Please reinstall it manually\n", green(pkgName))
		return nil
	}

	release, err := checkGithubTag(i.Source, "")
	if err != nil {
		return fmt.Errorf("failed to check version: %w", err)
	}

	if i.Version == release.TagName {
		fmt.Printf("%s is already up to date (%s)\n", green(pkgName), i.Version)
		return nil
	}

	pkg, err := resolveGithub(i.Source, i.PkgName, release)
	if err != nil {
		return err
	}

	return update(i, pkg)
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
		item    registryItem
		release *ghRelease
	}

	var updates []pending

	fmt.Println("Checking for updates, might take some time...")
	for _, i := range reg.Packages {
		if i.SourceType != "local" {
			release, err := fetchGhRelease(i.Source, "")
			if err != nil {
				fmt.Printf("Skipping %s: %v\n", i.PkgName, err)
				continue
			}

			if release.TagName != i.Version {
				updates = append(updates, pending{
					item:    i,
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
		fmt.Printf("  %s  %s -> %s\n", u.item.PkgName, u.item.Version, u.release.TagName)
	}

	if !confirm("\nUpdate all?") {
		return nil
	}

	for _, u := range updates {
		fmt.Println()
		pkg, err := resolveGithub(u.item.Source, u.item.PkgName, u.release)
		if err != nil {
			return err
		}
		if err := update(u.item, pkg); err != nil {
			fmt.Printf("failed to update %s: %v\n", u.item.PkgName, err)
		}
	}

	return nil
}
