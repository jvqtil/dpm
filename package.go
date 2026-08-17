package main

import (
	"fmt"
)

type Package struct {
	Name        string `json:"pkg_name"`
	CurrentTag  Tag    `json:"current_tag"`
	Tags        []Tag  `json:"tags"`
	SourceType  string `json:"source_type"`
	Source      string `json:"source"`
	BinaryPath  string `json:"binary_path"`
	InstalledAt string `json:"installed_at"`
	LastUpdated string `json:"last_updated"`
}

type Tag struct {
	Name      string `json:"tag_name"`
	AssetName string `json:"asset_name"`
	AssetURL  string `json:"asset_url"`
	AssetPath string `json:"asset_path"`
	AssetSize int64  `json:"asset_size"`
}

func (p *Package) AddTag(t Tag) {
	for i, existing := range p.Tags {
		if existing.Name == t.Name {
			p.Tags[i] = t
			return
		}
	}
	p.Tags = append(p.Tags, t)
}

func (p *Package) SetCurrentTag(t Tag) error {
	for _, existing := range p.Tags {
		if existing.Name == t.Name {
			p.CurrentTag = existing
			return nil
		}
	}
	return fmt.Errorf("tag %q not found for package %s", t.Name, p.Name)
}

func (p *Package) RemoveTag(t Tag) error {
	for i, existing := range p.Tags {
		if existing.Name == t.Name {
			p.Tags = append(p.Tags[:i], p.Tags[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("tag %q not found for package %s", t.Name, p.Name)
}
