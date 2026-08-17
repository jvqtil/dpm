package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dustin/go-humanize"
)

func clearCache(dir string) error {
	cacheSize := dirSize(dir)
	cacheSizeFmt := humanize.Bytes(uint64((cacheSize)))
	if cacheSize == 0 {
		fmt.Println("Cache directory is empty, nothing to clear")
		return nil
	}
	if !confirm(fmt.Sprintf("Cache size in %s: %s.\nClear cache?", filepath.Base(dir), red(cacheSizeFmt))) {
		return nil
	}

	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("failed to clear cache directory: %w", err)
	}

	fmt.Printf("=> Cleared %s of cache\n", red(cacheSizeFmt))
	return nil
}

func clearPkgCache(name, tag string) error {
	reg, err := loadRegistry()
	if err != nil {
		return err
	}

	pkg, ok := reg.Packages[name]
	if !ok {
		return fmt.Errorf("package %s is not installed", name)
	}

	return clearCache(filepath.Join(pkg.ResolveCachePath(), tag))
}
