package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

func installPkg(input, explicitName string) error {
	source, tag := resolveTag(input)

	if explicitName == "" {
		explicitName = resolvePkgName(source)
	}
	explicitName = strings.ToLower(explicitName)

	reg, err := loadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	ePkg, exists := reg.Packages[explicitName]
	if exists {
		var tagVerb string
		if ePkg.CurrentTag.TagName != "" {
			tagVerb = "(" + ePkg.CurrentTag.TagName + ")"
		}
		fmt.Printf("=> %s is already installed %s\n", green(ePkg.Name), tagVerb)
		if !confirm("Reinstall?") {
			fmt.Println("Aborted")
			return nil
		}
	}

	var pkg *Package
	if isLocalPkg(source) {
		pkg, err = resolveLocal(source, tag, explicitName)
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
			pkg, err = resolveGithub(source, explicitName, release)
			if err != nil {
				return err
			}
		} else {
			pkg, err = resolveDirect(source, explicitName, tag)
			if err != nil {
				return err
			}
		}
	}

	border := strings.Repeat("═", 40)
	fmt.Println(border)
	fmt.Printf("%s%s (%s)\n\n", green(pkg.Name), getTagVerb(pkg.CurrentTag.TagName), pkg.SourceType)
	assetNameFmt := pkg.CurrentTag.AssetName
	if pkg.SourceType == "local" {
		assetNameFmt = pkg.CurrentTag.AssetPath
	}
	fmt.Printf("↓ %s ", assetNameFmt)
	if pkg.CurrentTag.AssetSize != 0 {
		fmt.Printf("(%s)", humanize.Bytes(uint64(pkg.CurrentTag.AssetSize)))
	}
	fmt.Printf("\n→ %s\n", pkg.BinaryPath)
	fmt.Println(border)

	if !exists && !confirm("Install this package?") {
		fmt.Println("Aborted")
		return nil
	}

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

	if err := cpToDest(src, pkg.BinaryPath); err != nil {
		return err
	}

	if exists {
		ePkg.AddTag(pkg.CurrentTag)
		if err := ePkg.SetCurrentTag(pkg.CurrentTag); err != nil {
			return err
		}
		ePkg.SourceType = pkg.SourceType
		ePkg.Source = pkg.Source
		ePkg.BinaryPath = pkg.BinaryPath
		ePkg.LastUpdated = time.Now().Format(time.RFC3339)
		reg.Packages[explicitName] = ePkg
	} else {
		pkg.InstalledAt = time.Now().Format(time.RFC3339)
		pkg.LastUpdated = pkg.InstalledAt
		pkg.Tags = []Tag{pkg.CurrentTag}
		reg.Packages[explicitName] = *pkg
	}

	if err := saveRegistry(reg); err != nil {
		rmBin(pkg.BinaryPath)
		return fmt.Errorf("failed to save registry, rolled back: %w", err)
	}

	fmt.Printf("=> Installed package %s%s\n", green(pkg.Name), getTagVerb(pkg.CurrentTag.TagName))
	return nil
}
