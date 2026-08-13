package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func list() error {
	reg, err := loadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}
	if len(reg.Packages) == 0 {
		fmt.Println("No packages installed")
		return nil
	}

	fmt.Printf("=> Installed packages (%d)\n", len(reg.Packages))
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	for _, i := range reg.Packages {
		fmt.Fprintf(w, "%s\t%s\n", i.PkgName, i.Version)
	}
	w.Flush()

	return nil
}

func listTarget(pkgName string) error {
	reg, err := loadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}
	i, ok := reg.Packages[strings.ToLower(pkgName)]
	if !ok {
		return fmt.Errorf("package %s is not installed", pkgName)
	}

	fmt.Printf("%s%s\n", green(i.PkgName), getTagVerb(i.Version))
	fmt.Printf("Source: %s\n\n", i.Source)
	fmt.Printf("Binary: %s\n", i.Binary)
	if i.LastUpdated != i.InstalledAt {
		fmt.Printf("Last Updated: %s\n", i.LastUpdated)
	}
	fmt.Printf("Installed: %s\n", i.InstalledAt)

	return nil
}
