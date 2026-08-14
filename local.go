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

func resolveLocal(source, tag_, name string) (*pkg, error) {
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

	return &pkg{
		Name:       name,
		SourceType: "local",
		Source:     source,
		BinaryPath: filepath.Join(cfg.BinDir, name),
		CurrentTag: tag{
			TagName:   tag_,
			AssetName: filepath.Base(source),
			AssetPath: source,
			AssetSize: file.Size(),
		},
		Tags: []tag{},
	}, nil
}
