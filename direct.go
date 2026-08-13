package main

import (
	"path/filepath"
	"strings"
)

func resolveDirect(source, name, tag string) (*pkg, error) {
	assetURL := source
	// Prepend https:// if the URL has no scheme
	if !strings.HasPrefix(source, "http://") && !strings.HasPrefix(source, "https://") {
		assetURL = "https://" + source
	}

	return &pkg{
		Name:       name,
		Version:    tag,
		SourceType: "direct",
		Source:     normalizeSource(source),
		AssetName:  filepath.Base(source),
		AssetURL:   assetURL,
	}, nil
}
