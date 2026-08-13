package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

func install(input, explicitName string) error {
	source, tag := resolveTag(input)

	if explicitName == "" {
		explicitName = resolvePkgName(source)
	}
	explicitName = strings.ToLower(explicitName)

	reg, err := loadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	e, exists := reg.Packages[explicitName]
	if exists {
		fmt.Printf("=> %s is already installed (%s)\n", green(e.PkgName), e.Version)
		if !confirm("Reinstall?") {
			fmt.Println("Aborted")
			return nil
		}
	}

	var pkg *pkg
	var release *ghRelease
	if isLocalPkg(source) {
		pkg, err = resolveLocal(source, tag, explicitName)
		if err != nil {
			return fmt.Errorf("failed to resolve: %w", err)
		}
	} else {
		source, err = normalizeSource(source)
		if err != nil {
			return err
		}

		switch sourceDomain(source) {
		case "github.com":
			release, err = checkGithubTag(source, tag)
			if err != nil {
				return err
			}

			pkg, err = resolveGithub(source, explicitName, release)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported source: %s", source)
		}

	}

	dest := filepath.Join(cfg.BinDir, pkg.Name)
	border := strings.Repeat("═", 40)
	fmt.Println(border)
	fmt.Printf("%s%s (%s)\n\n", green(pkg.Name), getTagVerb(pkg.Version), pkg.SourceType)
	assetNameFmt := pkg.AssetName
	if pkg.SourceType == "local" {
		assetNameFmt = pkg.AssetURL
	}
	suffix := ""
	if isArchive(pkg.AssetName) {
		suffix = " — " + cyan("archive")
	}
	fmt.Printf("↓ %s%s (%s)\n", assetNameFmt, suffix, humanize.Bytes(uint64(pkg.AssetSize)))
	fmt.Printf("→ %s\n", dest)
	fmt.Println(border)

	if !exists && !confirm("Install this package?") {
		fmt.Println("Aborted")
		return nil
	}

	src := pkg.Source
	if pkg.SourceType != "local" {
		src, err = resolveBinary(pkg)
		if err != nil {
			return err
		}
		if isArchive(pkg.AssetName) {
			if info, err := os.Stat(src); err == nil {
				fmt.Printf("%s %s (%s)\n", bold("Extracted:"), filepath.Base(src), humanize.Bytes(uint64(info.Size())))
			}
		}
	}

	if err := cpToDest(src, dest); err != nil {
		return err
	}

	reg.Packages[pkg.Name] = registryItem{
		PkgName:     pkg.Name,
		Version:     pkg.Version,
		SourceType:  pkg.SourceType,
		Source:      pkg.Source,
		AssetName:   pkg.AssetName,
		AssetURL:    pkg.AssetURL,
		Binary:      dest,
		InstalledAt: time.Now().Format("02 Jan 06 15:04"),
		LastUpdated: time.Now().Format("02 Jan 06 15:04"),
	}

	if err := saveRegistry(reg); err != nil {
		os.Remove(dest)
		return fmt.Errorf("failed to save registry, rolled back: %w", err)
	}

	fmt.Printf("=> Installed package %s%s\n", green(pkg.Name), getTagVerb(tag))
	return nil
}
