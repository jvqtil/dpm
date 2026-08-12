package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

func downloadAsset(pkg *pkg) (string, error) {
	destDir := filepath.Join(cfg.CacheDir, pkg.Source, pkg.Version)

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	dest := filepath.Join(destDir, pkg.AssetName)

	if i, err := os.Stat(dest); err == nil {
		if i.Size() == pkg.AssetSize {
			fmt.Printf("Already downloaded: %s\n", dest)
			return dest, nil
		}
		fmt.Println("Found incomplete download. Re-downloading...")
	}

	resp, err := http.Get(pkg.AssetURL)
	if err != nil {
		return "", fmt.Errorf("download request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download returned: %s", resp.Status)
	}

	out, err := os.Create(dest)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer out.Close()

	lengh := pkg.AssetSize
	if lengh == 0 {
		lengh = resp.ContentLength
	}
	bar := progressbar.DefaultBytes(lengh, "Downloading: ")
	if _, err := io.Copy(io.MultiWriter(out, bar), resp.Body); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return dest, nil
}
