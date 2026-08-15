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

	if i, err := os.Stat(dest); err == nil && tag.AssetSize != 0 {
		if i.Size() == tag.AssetSize {
			return dest, nil
		}
	}

	if tag.AssetSize == 0 {
		if i, err := os.Stat(dest); err == nil && i.Size() > 0 {
			headResp, err := http.Head(tag.AssetURL)
			if err == nil {
				defer headResp.Body.Close()
				if headResp.StatusCode == http.StatusOK {
					if headResp.ContentLength > 0 && headResp.ContentLength == i.Size() {
						return dest, nil
					}
				}
			}
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
