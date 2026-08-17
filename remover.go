package main

import (
	"fmt"
	"strings"
)

func removePkg(pkgName string) error {
	reg, err := loadRegistry()
	if err != nil {
		return err
	}

	pkg, ok := reg.Packages[strings.ToLower(pkgName)]
	if !ok {
		return fmt.Errorf("package %s is not installed", pkgName)
	}

	if !confirm(fmt.Sprintf("Remove %s?", green(pkgName))) {
		fmt.Println("Aborted")
		return nil
	}

	if err := rmBin(pkg.BinaryPath); err != nil {
		return fmt.Errorf("failed to remove binary: %w", err)
	}

	delete(reg.Packages, strings.ToLower(pkgName))

	if err := saveRegistry(reg); err != nil {
		return fmt.Errorf("failed to update registry: %w", err)
	}

	fmt.Printf("=> Removed %s\n", green(pkgName))

	return nil
}
