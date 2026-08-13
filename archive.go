package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func isArchive(name string) bool {
	name = strings.ToLower(name)
	return strings.HasSuffix(name, ".tar.gz") ||
		strings.HasSuffix(name, ".tgz") ||
		strings.HasSuffix(name, ".zip")
}

func isDocFile(name string) bool {
	name = strings.ToLower(name)
	base := filepath.Base(name)
	docNames := []string{"license", "readme", "changelog", "notice"}
	for _, d := range docNames {
		if strings.HasPrefix(base, d) {
			return true
		}
	}
	return strings.HasSuffix(name, ".md") || strings.HasSuffix(name, ".txt")
}

func extractArchive(archivePath string) ([]string, error) {
	destDir := archivePath + "-extracted"
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp dir for extracted archive: %w", err)
	}

	name := strings.ToLower(archivePath)
	switch {
	case strings.HasSuffix(name, ".tar.gz"), strings.HasSuffix(name, ".tgz"):
		return extractTarGz(archivePath, destDir)
	case strings.HasSuffix(name, ".zip"):
		return extractZip(archivePath, destDir)
	default:
		return nil, fmt.Errorf("unsupported archive: %s", archivePath)
	}
}

func extractTarGz(archivePath, destDir string) ([]string, error) {
	f, err := os.Open(archivePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("failed to open gzip: %w", err)
	}
	defer gz.Close()

	var extracted []string
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("tar read error: %w", err)
		}
		if hdr.Typeflag != tar.TypeReg {
			continue
		}

		outPath := filepath.Join(destDir, filepath.Base(hdr.Name))
		out, err := os.Create(outPath)
		if err != nil {
			return nil, err
		}
		if _, err := io.Copy(out, tr); err != nil {
			out.Close()
			return nil, err
		}
		out.Close()
		os.Chmod(outPath, os.FileMode(hdr.Mode))
		extracted = append(extracted, outPath)
	}
	return extracted, nil
}

func extractZip(archivePath, destDir string) ([]string, error) {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var extracted []string
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return nil, err
		}

		outPath := filepath.Join(destDir, filepath.Base(f.Name))
		out, err := os.Create(outPath)
		if err != nil {
			rc.Close()
			return nil, err
		}
		if _, err := io.Copy(out, rc); err != nil {
			out.Close()
			rc.Close()
			return nil, err
		}
		out.Close()
		rc.Close()
		os.Chmod(outPath, os.FileMode(f.FileInfo().Mode()))
		extracted = append(extracted, outPath)
	}
	return extracted, nil
}
