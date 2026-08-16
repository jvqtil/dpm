package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/dustin/go-humanize"
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
		fmt.Fprintf(w, "%s\t%s\n", p.Name, p.CurrentTag.Name)
	}
	w.Flush()

	return nil
}

func showPkgInfo(pkgName string) error {
	reg, err := loadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}
	p, ok := reg.Packages[strings.ToLower(pkgName)]
	if !ok {
		return fmt.Errorf("package %s is not installed", pkgName)
	}

	fmt.Println(border)
	fmt.Printf("%s - %s\n", green(p.Name), p.CurrentTag.Name)
	fmt.Printf("%s\n", p.Source)
	fmt.Printf("%s\n", p.BinaryPath)
	fmt.Println()

	tags := make([]string, len(p.Tags))
	for i, t := range p.Tags {
		tags[i] = t.Name
	}
	sort.Strings(tags)
	slices.Reverse(tags)

	for _, name := range tags {
		if name == p.CurrentTag.Name {
			fmt.Printf(" %s - current\n", name)
		} else {
			fmt.Printf(" %s\n", name)
		}
	}

	fmt.Println()
	if p.LastUpdated != p.InstalledAt {
		t, _ := time.Parse(time.RFC3339, p.LastUpdated)
		fmt.Printf("last updated: %s\n", humanize.Time(t))
	}
	t, _ := time.Parse(time.RFC3339, p.InstalledAt)
	fmt.Printf("installed: %s\n", t.Format("02 Jan 06 15:04"))
	fmt.Println(border)

	return nil
}

func showPkgInfoJSON(pkgName string) error {
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
