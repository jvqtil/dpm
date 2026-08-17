package main

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

func showPkgInfo(pkgName string, jsonOut bool) error {
	reg, err := loadRegistry()
	if err != nil {
		return err
	}
	pkg, ok := reg.Packages[strings.ToLower(pkgName)]
	if !ok {
		return fmt.Errorf("package %s is not installed", pkgName)
	}

	if jsonOut {
		data, err := json.MarshalIndent(pkg, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	} else {
		fmt.Println(border)
		fmt.Printf("%s - %s\n", green(pkg.Name), pkg.CurrentTag.Name)
		fmt.Printf("%s\n", pkg.Source)
		fmt.Printf("%s\n", pkg.BinaryPath)
		fmt.Println()

		tags := make([]string, len(pkg.Tags))
		for i, t := range pkg.Tags {
			tags[i] = t.Name
		}
		slices.Sort(tags)
		slices.Reverse(tags)

		for _, name := range tags {
			if name == pkg.CurrentTag.Name {
				fmt.Printf(" %s - current\n", name)
			} else {
				fmt.Printf(" %s\n", name)
			}
		}

		fmt.Println()
		if pkg.LastUpdated != pkg.InstalledAt {
			t, _ := time.Parse(time.RFC3339, pkg.LastUpdated)
			fmt.Printf("last updated: %s\n", humanize.Time(t))
		}
		t, _ := time.Parse(time.RFC3339, pkg.InstalledAt)
		fmt.Printf("installed: %s\n", t.Format("02 Jan 06 15:04"))
		fmt.Println(border)
	}
	return nil
}
