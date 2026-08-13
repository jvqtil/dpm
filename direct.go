package main

import (
	"path/filepath"
	"strings"
)

func resolveDirect(source, name, tag string) (*pkg, error) {
	assetURL := source

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
