package main

import (
	"fmt"
	"strings"
	"time"
)

func updatePkg(pkg Package, newTag Tag) error {
	src, err := resolveAsset(pkg, &newTag)
	if err != nil {
		return err
	}
	if newTag.AssetPath == "" {
		newTag.AssetPath = src
	}

	dest := pkg.BinaryPath

	if err := cpToDest(src, dest); err != nil {
		return fmt.Errorf("copy failed: %w", err)
	}

	reg, err := loadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	oldTag := pkg.CurrentTag.Name
	pkg.AddTag(newTag)
	if err := pkg.SetCurrentTag(newTag); err != nil {
		return err
	}
	pkg.LastUpdated = time.Now().Format(time.RFC3339)
	reg.Packages[pkg.Name] = pkg

	if err := saveRegistry(reg); err != nil {
		return fmt.Errorf("failed to save registry: %w", err)
	}

	fmt.Printf("=> Updated package %s from %s to %s\n", green(pkg.Name), oldTag, newTag.Name)
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

	var newPkg *Package
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

		if existingPkg.CurrentTag.Name == release.Name {
			fmt.Printf("%s is already up to date (%s)\n", green(pkgName), existingPkg.CurrentTag.Name)
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
		item    Package
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

			if release.Name != p.CurrentTag.Name {
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
		fmt.Printf("  %s  %s -> %s\n", u.item.Name, u.item.CurrentTag.Name, u.release.Name)
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
