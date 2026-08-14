package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"
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
		t, _ := time.Parse(time.RFC3339, p.LastUpdated)
		fmt.Printf("Last Updated: %s\n", t.Format("02 Jan 06 15:04"))
	}
	t, _ := time.Parse(time.RFC3339, p.InstalledAt)
	fmt.Printf("Installed: %s\n", t.Format("02 Jan 06 15:04"))

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
