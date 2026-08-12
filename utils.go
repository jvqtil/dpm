package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
)

var (
	green = color.New(color.FgGreen).SprintFunc()
	red   = color.New(color.FgRed).SprintFunc()
)

func resolveTag(input string) (source, tag string) {
	s := strings.TrimSpace(input)
	if idx := strings.LastIndex(s, "@"); idx != -1 {
		return s[:idx], s[idx+1:]
	}
	return s, ""
}

func resolvePkgName(input string) string {
	s := strings.TrimSuffix(input, "/")
	s = strings.TrimSuffix(s, ".git")
	if idx := strings.LastIndex(s, "/"); idx != -1 {
		return s[idx+1:]
	}
	return s
}

func getTagVerb(tag string) string {
	var tagVerb string
	if tag == "local" {
		return ""
	}
	if tag != "" {
		tagVerb = "@" + tag
	}
	return tagVerb
}

func confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y/N]: ", prompt)

	line, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	answer := strings.ToLower(strings.TrimSpace(line))
	return answer == "y" || answer == "yes"
}

func runSudo(args ...string) error {
	cmd := exec.Command("sudo", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func cpToDest(src, dest string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	err = os.MkdirAll(filepath.Dir(dest), 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(dest, data, 0755)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		return fmt.Errorf("missing directory, create it yourself: %w", err)
	}
	if !os.IsPermission(err) {
		return err
	}

	fmt.Println("=> Copying binary, may require password")

	if err := runSudo("mkdir", "-p", filepath.Dir(dest)); err != nil {
		return fmt.Errorf("failed to create a directory: %w", err)
	}

	if err := runSudo("cp", src, dest); err != nil {
		return fmt.Errorf("copying failed: %w", err)
	}

	if err := runSudo("chmod", "0755", dest); err != nil {
		return fmt.Errorf("chmod failed: %w", err)
	}

	return nil
}

func rmBin(target string) error {
	err := os.Remove(target)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		fmt.Printf("Note: %s was already missing, cleaning up registry item\n", target)
		return nil
	}
	if !os.IsPermission(err) {
		return err
	}

	if err := runSudo("rm", target); err != nil {
		return fmt.Errorf("sudo rm failed: %w", err)
	}

	return nil
}

func dirSize(path string) (size int64) {
	filepath.Walk(path, func(_ string, info os.FileInfo, _ error) error {
		if info != nil && !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return
}

func clearCache() error {
	cacheSize := dirSize(cfg.CacheDir)
	cacheSizeFmt := humanize.Bytes(uint64((cacheSize)))
	if cacheSize == 0 {
		fmt.Println("=> Cache directory is empty, nothing to clear")
		return nil
	}
	if !confirm(fmt.Sprintf("=> Cache directory size: %s.\nClear cache?", red(cacheSizeFmt))) {
		return nil
	}

	if err := os.RemoveAll(cfg.CacheDir); err != nil {
		return fmt.Errorf("failed to clear cache directory: %w", err)
	}

	fmt.Printf("=> Cleared %s of cache\n", red(cacheSizeFmt))
	return nil
}
