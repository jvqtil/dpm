package main

import (
	"fmt"
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
