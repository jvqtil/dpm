package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

func downloadAsset(tag *Tag, destDir string) (string, error) {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	dest := filepath.Join(destDir, tag.AssetName)

	// For packages with known size, check if cached file matches
	if i, err := os.Stat(dest); err == nil && tag.AssetSize != 0 {
		if i.Size() == tag.AssetSize {
			fmt.Printf("=> Already downloaded: %s\n", dest)
			return dest, nil
		}
		fmt.Println("=> Found incomplete download. Re-downloading...")
	}

	// For direct packages (AssetSize == 0), perform conditional request to validate cache
	if tag.AssetSize == 0 {
		if i, err := os.Stat(dest); err == nil && i.Size() > 0 {
			// Perform HEAD request to check if cached file is still valid
			headResp, err := http.Head(tag.AssetURL)
			if err == nil {
				defer headResp.Body.Close()
				if headResp.StatusCode == http.StatusOK {
					// Check if content length matches cached file
					if headResp.ContentLength > 0 && headResp.ContentLength == i.Size() {
						fmt.Printf("=> Already downloaded: %s\n", dest)
						return dest, nil
					}
					// If ETag is available and file exists, could add conditional GET here
					// For now, re-download if size doesn't match
				}
			}
			fmt.Println("=> Found cached file, validating...")
		}
	}

	resp, err := http.Get(tag.AssetURL)
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

	lengh := tag.AssetSize
	if lengh == 0 {
		lengh = resp.ContentLength
	}
	bar := progressbar.DefaultBytes(lengh, "Downloading: ")
	if _, err := io.Copy(io.MultiWriter(out, bar), resp.Body); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return dest, nil
}
