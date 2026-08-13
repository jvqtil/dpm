package main

import "path/filepath"

func resolveDirect(source, name, tag string) (*pkg, error) {
	return &pkg{
		Name:       name,
		Version:    tag,
		SourceType: "direct",
		Source:     normalizeSource(source),
		AssetName:  filepath.Base(source),
		AssetURL:   source,
	}, nil
}
