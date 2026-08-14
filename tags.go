package main

import "fmt"

func (p *Package) AddTag(t Tag) {
	for i, existing := range p.Tags {
		if existing.TagName == t.TagName {
			p.Tags[i] = t
			return
		}
	}
	p.Tags = append(p.Tags, t)
}

func (p *Package) SetCurrentTag(t Tag) error {
	for _, existing := range p.Tags {
		if existing.TagName == t.TagName {
			p.CurrentTag = existing
			return nil
		}
	}
	return fmt.Errorf("tag %q not found for package %s", t.TagName, p.Name)
}
