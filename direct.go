package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func resolveDirect(source, name, tag_ string) (*pkg, error) {
	if strings.Contains(name, "/") || strings.Contains(name, "\\") || strings.HasPrefix(name, "..") {
		return nil, fmt.Errorf("invalid package name: %s contains path separators or traversal components", name)
	}

	assetURL := source

	if !strings.HasPrefix(source, "http://") && !strings.HasPrefix(source, "https://") {
		assetURL = "https://" + source
	}

	if strings.Contains(name, "/") || strings.Contains(name, "\\") || strings.HasPrefix(name, "..") {
		return nil, fmt.Errorf("invalid package name: %s contains path separators or traversal components", name)
	}

	return &pkg{
		Name:       name,
		SourceType: "direct",
		Source:     normalizeSource(source),
		BinaryPath: filepath.Join(cfg.BinDir, name),
		CurrentTag: tag{
			TagName:   tag_,
			AssetName: filepath.Base(source),
			AssetURL:  assetURL,
		},
		Tags: []tag{},
	}, nil
}
