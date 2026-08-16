package main

import (
	"fmt"
	"strings"
	"time"
)

func updatePkg(pkg Package, tag Tag) error {
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

	oldTag := pkg.CurrentTag.Name
	pkg.AddTag(tag)
	if err := pkg.SetCurrentTag(tag); err != nil {
		return err
	}
	pkg.LastUpdated = time.Now().Format(time.RFC3339)
	reg.Packages[pkg.Name] = pkg

	if err := saveRegistry(reg); err != nil {
		return err
	}

	fmt.Printf("=> Updated package %s from %s to %s\n", green(pkg.Name), oldTag, tag.Name)
	return nil
}

func updateTarget(pkgName string) error {
	reg, err := loadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	ePkg, ok := reg.Packages[strings.ToLower(pkgName)]
	if !ok {
		return fmt.Errorf("package %s is not installed", pkgName)
	}

	var newPkg *Package
	switch ePkg.SourceType {
	case "local":
		fmt.Printf("=> %s is a local package. Local packages can't be updated. Please reinstall it manually\n", green(pkgName))
		return nil
	case "direct":
		fmt.Printf("=> %s was installed from direct URL. dpm can't fetch updates for it. Please reinstall it manually\n", green(pkgName))
		return nil
	case "github.com":
		release, err := checkGithubTag(ePkg.Source, "")
		if err != nil {
			return fmt.Errorf("failed to check version: %w", err)
		}

		if ePkg.CurrentTag.Name == release.Name {
			fmt.Printf("%s is already up to date (%s)\n", green(pkgName), ePkg.CurrentTag.Name)
			return nil
		}

		newPkg, err = resolveGithub(ePkg.Source, ePkg.Name, release)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported source type %s", ePkg.SourceType)
	}

	return switchTag(ePkg, newPkg.CurrentTag, fmt.Sprintf("Updated package %s to %s", green(ePkg.Name), newPkg.CurrentTag.Name))
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
