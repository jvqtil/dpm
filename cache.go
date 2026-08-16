package main

import (
	"fmt"
	"os"

	"github.com/dustin/go-humanize"
)

func clearCache() error {
	cacheSize := dirSize(cfg.CacheDir)
	cacheSizeFmt := humanize.Bytes(uint64((cacheSize)))
	if cacheSize == 0 {
		fmt.Println("Cache directory is empty, nothing to clear")
		return nil
	}
	if !confirm(fmt.Sprintf("Cache directory size: %s.\nClear cache?", red(cacheSizeFmt))) {
		return nil
	}

	if err := os.RemoveAll(cfg.CacheDir); err != nil {
		return fmt.Errorf("failed to clear cache directory: %w", err)
	}

	fmt.Printf("=> Cleared %s of cache\n", red(cacheSizeFmt))
	return nil
}
