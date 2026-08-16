package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func resolveAsset(pkg Package, tag *Tag) (string, error) {
	src, err := downloadAsset(tag, resolveCachePath(pkg, *tag))
	if err != nil {
		return "", err
	}

	if tag.AssetSize == 0 {
		if info, err := os.Stat(src); err == nil {
			tag.AssetSize = info.Size()
		}
	}

	kind, err := sniffArchiveKind(src)
	if err != nil {
		return "", fmt.Errorf("failed to sniff archive kind: %w", err)
	}
	if kind == NotArchive {
		return src, nil
	}

	files, err := extractArchive(src, kind)
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

	picked, err := picker(names, fmt.Sprintf("Found %d files in archive - host: %s/%s", len(names), runtime.GOOS, runtime.GOARCH))
	if err != nil {
		return "", err
	}

	return nameToPath[picked], nil
}

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

func resolveCachePath(p Package, t Tag) string {
	var cachePath string
	switch p.SourceType {
	case "github.com":
		cachePath = filepath.Join(cfg.CacheDir, p.Source, t.Name)
	default:
		cachePath = filepath.Join(cfg.CacheDir, getSourceDomain(normalizeSource(p.Source)), p.Name, t.Name)
	}
	return cachePath
}
