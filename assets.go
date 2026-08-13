package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func matchByOsArch(names []string) (string, bool) {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	osAliases := map[string][]string{
		"darwin":  {"darwin", "macos", "mac", "osx"},
		"linux":   {"linux"},
		"windows": {"windows", "win"},
	}

	archAliases := map[string][]string{
		"amd64": {"amd64", "x86_64", "x64"},
		"arm64": {"arm64", "aarch64"},
	}

	skipSfx := []string{".sha256", ".sig", ".asc", ".txt", ".sbom"}

	for _, n := range names {
		name := strings.ToLower(n)
		skip := false
		for _, sfx := range skipSfx {
			if strings.HasSuffix(name, sfx) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		var osMatch, archMatch bool
		for _, x := range osAliases[goos] {
			osMatch = osMatch || strings.Contains(name, x)
		}
		for _, x := range archAliases[goarch] {
			archMatch = archMatch || strings.Contains(name, x)
		}

		if osMatch && archMatch {
			return n, true
		}
	}

	return "", false
}

func pickAsset(names []string, prompt string) (string, error) {
	if len(names) == 0 {
		return "", fmt.Errorf("nothing to pick from")
	}

	fmt.Printf("\n%s\n", prompt)
	for i, n := range names {
		fmt.Printf("%d) %s\n", i+1, n)
	}

	fmt.Printf("\nYour pick? ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}
	line = strings.TrimSpace(line)
	if line == "" {
		os.Exit(0)
	}

	choice, err := strconv.Atoi(line)
	if err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	if choice < 1 || choice > len(names) {
		return "", fmt.Errorf("choice %d out of range", choice)
	}

	return names[choice-1], nil
}

func resolveBinary(pkg *pkg) (string, error) {
	var destDir string
	switch pkg.SourceType {
	case "github.com":
		destDir = filepath.Join(cfg.CacheDir, pkg.Source, pkg.Version)
	default:
		destDir = filepath.Join(cfg.CacheDir, sourceDomain(pkg.Source), pkg.Name, pkg.Version)
	}
	src, err := downloadAsset(pkg, destDir)
	if err != nil {
		return "", err
	}

	if !isArchive(pkg.AssetName) {
		return src, nil
	}

	files, err := extractArchive(src)
	if err != nil {
		return "", fmt.Errorf("failed to extract archive: %w", err)
	}

	var candidates []string
	for _, f := range files {
		if !isDocFile(f) {
			candidates = append(candidates, f)
		}
	}
	if len(candidates) == 0 {
		return "", fmt.Errorf("archive contains no binaries")
	}
	if len(candidates) == 1 {
		return candidates[0], nil
	}

	for _, c := range candidates {
		base := filepath.Base(c)
		if strings.EqualFold(base, pkg.Name) || strings.EqualFold(strings.TrimSuffix(base, filepath.Ext(base)), pkg.Name) {
			return c, nil
		}
	}

	names := make([]string, len(candidates))
	nameToPath := make(map[string]string, len(candidates))
	for i, c := range candidates {
		base := filepath.Base(c)
		names[i] = base
		nameToPath[base] = c
	}
	if matched, ok := matchByOsArch(names); ok {
		return nameToPath[matched], nil
	}

	picked, err := pickAsset(names, fmt.Sprintf("Found %d files in archive - host: %s/%s", len(names), runtime.GOOS, runtime.GOARCH))
	if err != nil {
		return "", err
	}
	return nameToPath[picked], nil
}
