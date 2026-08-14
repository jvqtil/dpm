package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func listAll() error {
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
	for _, p := range reg.Packages {
		fmt.Fprintf(w, "%s\t%s\n", p.Name, p.CurrentTag.TagName)
	}
	w.Flush()

	return nil
}

func listTarget(pkgName string) error {
	reg, err := loadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}
	p, ok := reg.Packages[strings.ToLower(pkgName)]
	if !ok {
		return fmt.Errorf("package %s is not installed", pkgName)
	}

	fmt.Printf("%s%s\n", green(p.Name), getTagVerb(p.CurrentTag.TagName))
	fmt.Printf("Source: %s\n\n", p.Source)
	fmt.Printf("Binary: %s\n", p.BinaryPath)
	if p.LastUpdated != p.InstalledAt {
		fmt.Printf("Last Updated: %s\n", p.LastUpdated)
	}
	fmt.Printf("Installed: %s\n", p.InstalledAt)

	return nil
}

func listTargetJSON(pkgName string) error {
	reg, err := loadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}
	pkg, ok := reg.Packages[strings.ToLower(pkgName)]
	if !ok {
		return fmt.Errorf("package %s is not installed", pkgName)
	}

	data, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	return nil
}
