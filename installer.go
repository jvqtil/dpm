package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

func installPkg(input, explicitName string) error {
	source, tag_ := resolveTag(input)

	if explicitName == "" {
		explicitName = resolvePkgName(source)
	}
	explicitName = strings.ToLower(explicitName)

	reg, err := loadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	existingPkg, exists := reg.Packages[explicitName]
	if exists {
		var tagVerb string
		if existingPkg.CurrentTag.TagName != "" {
			tagVerb = "(" + existingPkg.CurrentTag.TagName + ")"
		}
		fmt.Printf("=> %s is already installed %s\n", green(existingPkg.Name), tagVerb)
		if !confirm("Reinstall?") {
			fmt.Println("Aborted")
			return nil
		}
	}

	var pkg *pkg
	if isLocalPkg(source) {
		pkg, err = resolveLocal(source, tag_, explicitName)
		if err != nil {
			return fmt.Errorf("failed to resolve: %w", err)
		}
	} else {
		if isGithubSource(normalizeSource(source)) {
			source = normalizeSource(source)
			release, err := checkGithubTag(source, tag_)
			if err != nil {
				return err
			}
			pkg, err = resolveGithub(source, explicitName, release)
			if err != nil {
				return err
			}
		} else {
			pkg, err = resolveDirect(source, explicitName, tag_)
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
	suffix := ""
	if isArchive(pkg.CurrentTag.AssetName) {
		suffix = " — " + cyan("archive")
	}
	fmt.Printf("↓ %s%s ", assetNameFmt, suffix)
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
		src, err = resolveAsset(pkg)
		if err != nil {
			return err
		}
		if pkg.CurrentTag.AssetPath == "" {
			pkg.CurrentTag.AssetPath = src
		}
		if isArchive(pkg.CurrentTag.AssetName) {
			if info, err := os.Stat(src); err == nil {
				fmt.Printf("%s %s (%s)\n", bold("Extracted:"), filepath.Base(pkg.CurrentTag.AssetPath), humanize.Bytes(uint64(info.Size())))
			}
		}
	}

	if err := cpToDest(src, pkg.BinaryPath); err != nil {
		return err
	}

	if exists {
		found := -1
		for i, t := range existingPkg.Tags {
			if t.TagName == pkg.CurrentTag.TagName {
				found = i
				break
			}
		}
		if found != -1 {
			existingPkg.Tags[found].AssetURL = pkg.CurrentTag.AssetURL
			existingPkg.Tags[found].AssetPath = pkg.CurrentTag.AssetPath
			existingPkg.Tags[found].AssetSize = pkg.CurrentTag.AssetSize
		} else {
			existingPkg.Tags = append(existingPkg.Tags, pkg.CurrentTag)
		}
		existingPkg.CurrentTag = pkg.CurrentTag
		existingPkg.LastUpdated = time.Now().Format("02 Jan 06 15:04")
		reg.Packages[explicitName] = existingPkg
	} else {
		pkg.InstalledAt = time.Now().Format("02 Jan 06 15:04")
		pkg.LastUpdated = pkg.InstalledAt
		pkg.Tags = []tag{pkg.CurrentTag}
		reg.Packages[explicitName] = *pkg
	}

	if err := saveRegistry(reg); err != nil {
		rmBin(pkg.BinaryPath)
		return fmt.Errorf("failed to save registry, rolled back: %w", err)
	}

	fmt.Printf("=> Installed package %s%s\n", green(pkg.Name), getTagVerb(pkg.CurrentTag.TagName))
	return nil
}
