package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type ghRelease struct {
	TagName string    `json:"tag_name"`
	Assets  []ghAsset `json:"assets"`
}

type ghAsset struct {
	Name string `json:"name"`
	URL  string `json:"browser_download_url"`
	Size int64  `json:"size"`
}

func fetchGhRelease(source, tag string) (*ghRelease, error) {
	parts := strings.Split(source, "/")
	if len(parts) != 3 || !strings.EqualFold(parts[0], "github.com") {
		return nil, fmt.Errorf("not a github.com source: %s", source)
	}

	var url string
	if tag == "" {
		url = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", parts[1], parts[2])
	} else {
		url = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", parts[1], parts[2], tag)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		if tag == "" {
			return nil, fmt.Errorf("no releases found")
		} else {
			return nil, fmt.Errorf("tag %q not found", tag)
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api returned: %s", resp.Status)
	}

	var release ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}

	return &release, nil
}

func checkGithubTag(source, tag string) (*ghRelease, error) {
	fmt.Printf("=> Fetching releases of %s%s\n", source, getTagVerb(tag))
	release, err := fetchGhRelease(source, tag)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release: %w", err)
	}
	return release, nil
}

func resolveGithub(source, name string, release *ghRelease) (*pkg, error) {
	asset := matchGhAsset(release.Assets)
	if asset == nil {
		var err error
		asset, err = pickGhAsset(release.Assets)
		if err != nil {
			return nil, err
		}
	}

	return &pkg{
		Name:       name,
		Version:    release.TagName,
		SourceType: sourceDomain(source),
		Source:     source,
		AssetName:  asset.Name,
		AssetURL:   asset.URL,
		AssetSize:  asset.Size,
	}, nil
}

func matchGhAsset(assets []ghAsset) *ghAsset {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	osAliases := map[string][]string{
		"darwin":  {"darwin", "macos", "mac", "osx"},
		"linux":   {"linux"},
		"windows": {"windows", "win"},
	}

	archAliases := map[string][]string{
		"amd64": {"amd64", "x86_64", "x64"},
		"arm64": {"arm64", "aarch64"},
	}

	skipSfx := []string{".sha256", ".sig", ".asc", ".txt", ".sbom"}

	for _, a := range assets {
		name := strings.ToLower(a.Name)
		skip := false
		for _, sfx := range skipSfx {
			if strings.HasSuffix(name, sfx) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		var osMatch, archMatch bool
		for _, x := range osAliases[goos] {
			osMatch = osMatch || strings.Contains(name, x)
		}
		for _, x := range archAliases[goarch] {
			archMatch = archMatch || strings.Contains(name, x)
		}

		if osMatch && archMatch {
			return &a
		}
	}

	return nil
}

func pickGhAsset(assets []ghAsset) (*ghAsset, error) {
	if len(assets) == 0 {
		return nil, fmt.Errorf("no assets found in latest release")
	}

	fmt.Printf("\nFound %d assets - host: %s/%s\n", len(assets), runtime.GOOS, runtime.GOARCH)
	for i, a := range assets {
		fmt.Printf("%d) %s\n", i+1, a.Name)
	}

	fmt.Printf("\nYour pick? ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}
	line = strings.TrimSpace(line)
	if line == "" {
		os.Exit(0)
	}

	choice, err := strconv.Atoi(line)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}
	if choice < 1 || choice > len(assets) {
		return nil, fmt.Errorf("choice %d out of range", choice)
	}

	return &assets[choice-1], nil
}
