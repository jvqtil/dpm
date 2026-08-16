package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var (
	green = color.New(color.FgGreen).SprintFunc()
	red   = color.New(color.FgRed).SprintFunc()
	cyan  = color.New(color.FgCyan).SprintFunc()
	bold  = color.New(color.Bold).SprintFunc()

	border = strings.Repeat("═", 40)
)

func resolvePkgName(input string) string {
	s := strings.TrimSuffix(input, "/")
	s = strings.TrimSuffix(s, ".git")
	if idx := strings.LastIndex(s, "/"); idx != -1 {
		return s[idx+1:]
	}
	return s
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

func picker(items []string, prompt string) (string, error) {
	if len(items) == 0 {
		return "", fmt.Errorf("nothing to pick from")
	}

	fmt.Printf("\n%s\n", prompt)
	for i, n := range items {
		fmt.Printf("%d) %s\n", i+1, n)
	}

	fmt.Printf("\nYour pick? ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}
	line = strings.TrimSpace(line)
	if line == "" {
		os.Exit(0)
	}

	choice, err := strconv.Atoi(line)
	if err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}
	if choice < 1 || choice > len(items) {
		return "", fmt.Errorf("choice %d out of range", choice)
	}

	return items[choice-1], nil
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
