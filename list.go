package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func listAll() error {
	reg, err := loadRegistry()
	if err != nil {
		return err
	}
	if len(reg.Packages) == 0 {
		fmt.Println("No packages installed")
		return nil
	}

	fmt.Printf("Installed packages (%d)\n", len(reg.Packages))
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', 0)
	for _, p := range reg.Packages {
		fmt.Fprintf(w, "%s\t%s\n", p.Name, p.CurrentTag.Name)
	}
	w.Flush()

	return nil
}
