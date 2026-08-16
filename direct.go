package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func resolveDirect(source, name, tag string) (*Package, error) {
	if strings.Contains(name, "/") || strings.Contains(name, "\\") || strings.HasPrefix(name, "..") {
		return nil, fmt.Errorf("invalid package name: %s contains path separators or traversal components", name)
	}
	if !strings.Contains(source, ".") {
		return nil, fmt.Errorf("'%s' doesnt look like a valid source", source)
	}

	assetURL := source

	if !strings.HasPrefix(source, "http://") && !strings.HasPrefix(source, "https://") {
		assetURL = "https://" + source
	}

	return &Package{
		Name:       name,
		SourceType: "direct",
		Source:     normalizeSource(source),
		BinaryPath: filepath.Join(cfg.BinDir, name),
		CurrentTag: Tag{
			Name:      tag,
			AssetName: filepath.Base(source),
			AssetURL:  assetURL,
		},
		Tags: []Tag{},
	}, nil
}
