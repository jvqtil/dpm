package main

import (
	"fmt"
	"strings"
	"time"
)

func switchTag(pkg Package, tag Tag, msg string) error {
	src, err := resolveAsset(pkg, &tag)
	if err != nil {
		return err
	}

	if err := cpToDest(src, pkg.BinaryPath); err != nil {
		return err
	}

	reg, err := loadRegistry()
	if err != nil {
		return err
	}

	pkg.AddTag(tag)
	if err := pkg.SetCurrentTag(tag); err != nil {
		return err
	}
	pkg.LastUpdated = time.Now().Format(time.RFC3339)
	reg.Packages[pkg.Name] = pkg

	if err := saveRegistry(reg); err != nil {
		return err
	}

	fmt.Printf("=> %s\n", msg)
	return nil
}

func useTag(pkgName, tag string) error {
	reg, err := loadRegistry()
	if err != nil {
		return err
	}

	pkg, ok := reg.Packages[strings.ToLower(pkgName)]
	if !ok {
		return fmt.Errorf("package %s is not installed", pkgName)
	}

	if pkg.CurrentTag.Name == tag {
		fmt.Printf("%s is already on %s\n", green(pkg.Name), tag)
	}

	var target *Tag
	for i, t := range pkg.Tags {
		if t.Name == tag {
			target = &pkg.Tags[i]
			break
		}
	}
	if target == nil {
		return fmt.Errorf("tag %q not found for package %s", tag, pkg.Name)
	}

	return switchTag(pkg, *target, fmt.Sprintf("Switched %s to %s", green(pkg.Name), tag))
}
