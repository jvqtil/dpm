package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
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

func isGithubSource(source string) bool {
	parts := strings.Split(strings.Trim(source, "/"), "/")
	return len(parts) == 3 && strings.EqualFold(parts[0], "github.com")
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
	if strings.Contains(name, "/") || strings.Contains(name, "\\") || strings.HasPrefix(name, "..") {
		return nil, fmt.Errorf("invalid package name: %s contains path separators or traversal components", name)
	}

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
		SourceType: getSourceDomain(source),
		Source:     source,
		BinaryPath: filepath.Join(cfg.BinDir, name),
		CurrentTag: tag{
			TagName:   release.TagName,
			AssetName: asset.Name,
			AssetURL:  asset.URL,
			AssetSize: asset.Size,
		},
		Tags: []tag{},
	}, nil
}

func matchGhAsset(assets []ghAsset) *ghAsset {
	names := make([]string, len(assets))
	for i, a := range assets {
		names[i] = a.Name
	}

	matched, ok := matchByOsArch(names)
	if !ok {
		return nil
	}

	for i, a := range assets {
		if a.Name == matched {
			return &assets[i]
		}
	}

	return nil
}

func pickGhAsset(assets []ghAsset) (*ghAsset, error) {
	names := make([]string, len(assets))
	for i, a := range assets {
		names[i] = a.Name
	}

	picked, err := pickAsset(names, fmt.Sprintf("Found %d assets - host: %s/%s", len(assets), runtime.GOOS, runtime.GOARCH))
	if err != nil {
		return nil, err
	}

	for i, a := range assets {
		if a.Name == picked {
			return &assets[i], nil
		}
	}
	return nil, fmt.Errorf("internal error: pick not found")
}
