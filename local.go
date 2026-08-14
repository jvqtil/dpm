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

func resolveLocal(source, tag, name string) (*Package, error) {
	if strings.Contains(name, "/") || strings.Contains(name, "\\") || strings.HasPrefix(name, "..") {
		return nil, fmt.Errorf("invalid package name: %s contains path separators or traversal components", name)
	}

	if strings.HasPrefix(source, "~/") {
		homeDir, _ := os.UserHomeDir()
		source = filepath.Join(homeDir, source[2:])
	}

	source, _ = filepath.Abs(source)

	file, err := os.Stat(source)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	return &Package{
		Name:       name,
		SourceType: "local",
		Source:     source,
		BinaryPath: filepath.Join(cfg.BinDir, name),
		CurrentTag: Tag{
			TagName:   tag,
			AssetName: filepath.Base(source),
			AssetPath: source,
			AssetSize: file.Size(),
		},
		Tags: []Tag{},
	}, nil
}
