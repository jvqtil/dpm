package main

import (
	"fmt"
	"strings"
)

func remove(pkgName string) error {
	reg, err := loadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	i, ok := reg.Packages[strings.ToLower(pkgName)]
	if !ok {
		return fmt.Errorf("package %s is not installed", pkgName)
	}

	if !confirm(fmt.Sprintf("Remove %s (%s)?", green(pkgName), i.Version)) {
		fmt.Println("Aborted")
		return nil
	}

	if err := rmBin(i.Binary); err != nil {
		return fmt.Errorf("failed to remove binary: %w", err)
	}

	delete(reg.Packages, strings.ToLower(pkgName))

	if err := saveRegistry(reg); err != nil {
		return fmt.Errorf("failed to update registry: %w", err)
	}

	fmt.Printf("=> Removed %s\n", green(pkgName))

	return nil
}
