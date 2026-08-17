package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

func installPkg(source, tag, name string) error {
	if name == "" {
		name = resolvePkgName(source)
	}
	name = strings.ToLower(name)

	reg, err := loadRegistry()
	if err != nil {
		return err
	}

	ePkg, exists := reg.Packages[name]
	if exists {
		fmt.Printf("%s is already installed (%s)\n", green(ePkg.Name), ePkg.CurrentTag.Name)
		if !confirm("Reinstall?") {
			fmt.Println("Aborted")
			return nil
		}
	}

	// resolve package
	var pkg *Package
	if isLocalPkg(source) {
		pkg, err = resolveLocal(source, tag, name)
		if err != nil {
			return fmt.Errorf("failed to resolve: %w", err)
		}
	} else {
		if isGithubSource(normalizeSource(source)) {
			source = normalizeSource(source)
			release, err := checkGithubTag(source, tag)
			if err != nil {
				return err
			}
			pkg, err = resolveGithub(source, name, release)
			if err != nil {
				return err
			}
		} else {
			pkg, err = resolveDirect(source, name, tag)
			if err != nil {
				return err
			}
		}
	}

	// print package info & ask to install
	fmt.Println(border)
	fmt.Printf("%s - %s (%s)\n\n", green(pkg.Name), pkg.CurrentTag.Name, pkg.SourceType)
	assetNameFmt := pkg.CurrentTag.AssetName
	if pkg.SourceType == "local" {
		assetNameFmt = filepath.Base(pkg.CurrentTag.AssetPath)
	}
	fmt.Printf("↓ %s ", assetNameFmt)
	if pkg.CurrentTag.AssetSize != 0 {
		fmt.Printf("(%s)", humanize.Bytes(uint64(pkg.CurrentTag.AssetSize)))
	}
	fmt.Printf("\n→ %s\n", pkg.BinaryPath)
	fmt.Println(border)

	if !confirm("Install this package?") {
		fmt.Println("Aborted")
		return nil
	}

	// get source file (local/cache)
	var src string
	if pkg.SourceType == "local" {
		src = pkg.CurrentTag.AssetPath
	} else {
		src, err = resolveAsset(*pkg, &pkg.CurrentTag)
		if err != nil {
			return err
		}
		if pkg.CurrentTag.AssetPath == "" {
			pkg.CurrentTag.AssetPath = src
		}
	}

	// move binary
	if err := cpToDest(src, pkg.BinaryPath); err != nil {
		return err
	}

	// fill missing entries in package
	pkg.InstalledAt = time.Now().Format(time.RFC3339)
	pkg.LastUpdated = pkg.InstalledAt
	pkg.Tags = []Tag{pkg.CurrentTag}
	reg.Packages[name] = *pkg

	if err := saveRegistry(reg); err != nil {
		rmBin(pkg.BinaryPath)
		return fmt.Errorf("failed to save registry, rolled back: %w", err)
	}

	fmt.Printf("=> Installed package %s - %s\n", green(pkg.Name), pkg.CurrentTag.Name)
	return nil
}
