package main

import (
	"fmt"
	"os"
	"strings"
)

func removePkg(pkgName string) error {
	reg, err := loadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	p, ok := reg.Packages[strings.ToLower(pkgName)]
	if !ok {
		return fmt.Errorf("package %s is not installed", pkgName)
	}

	var tagVerb string
	if p.CurrentTag.TagName != "" {
		tagVerb = "(" + p.CurrentTag.TagName + ")"
	}
	if !confirm(fmt.Sprintf("Remove %s %s?", green(pkgName), tagVerb)) {
		fmt.Println("Aborted")
		return nil
	}

	// Back up the binary before removal
	var backupData []byte
	if data, err := os.ReadFile(p.BinaryPath); err == nil {
		backupData = data
	}

	if err := rmBin(p.BinaryPath); err != nil {
		return fmt.Errorf("failed to remove binary: %w", err)
	}

	delete(reg.Packages, strings.ToLower(pkgName))

	if err := saveRegistry(reg); err != nil {
		// Restore the binary if registry save fails
		if backupData != nil {
			os.WriteFile(p.BinaryPath, backupData, 0755)
		}
		return fmt.Errorf("failed to update registry: %w", err)
	}

	fmt.Printf("=> Removed %s\n", green(pkgName))

	return nil
}
