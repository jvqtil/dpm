package main

import (
	"path/filepath"
	"strings"
)

func resolveDirect(source, name, tag_ string) (*pkg, error) {
	assetURL := source

	if !strings.HasPrefix(source, "http://") && !strings.HasPrefix(source, "https://") {
		assetURL = "https://" + source
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
