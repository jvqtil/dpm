package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func isLocalPkg(input string) bool {
	return strings.HasPrefix(input, "./") || strings.HasPrefix(input, "/") || strings.HasPrefix(input, "~/")
}

func resolveLocal(source, tag, name string) (*pkg, error) {
	if strings.HasPrefix(source, "~/") {
		homeDir, _ := os.UserHomeDir()
		source = filepath.Join(homeDir, source[2:])
	}

	source, _ = filepath.Abs(source)

	if _, err := os.Stat(source); err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	if tag == "" {
		tag = "local"
	}

	return &pkg{
		Name:       name,
		Version:    tag,
		SourceType: "local",
		Source:     source,
		AssetName:  filepath.Base(source),
		AssetURL:   source,
	}, nil
}
