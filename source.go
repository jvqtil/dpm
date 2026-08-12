package main

import (
	"fmt"
	"strings"
)

func normalizeSource(input string) (string, error) {
	s := strings.TrimSpace(input)
	s = strings.TrimPrefix(s, "https://")
	s = strings.TrimPrefix(s, "http://")
	s = strings.TrimPrefix(s, "www.")
	s = strings.TrimSuffix(s, "/")
	s = strings.TrimSuffix(s, ".git")

	parts := strings.Split(s, "/")

	switch len(parts) {
	case 2:
		return fmt.Sprintf("%s/%s", cfg.AssumeSourceType, s), nil
	case 3:
		return s, nil
	default:
		return "", fmt.Errorf("invalid source: %s", input)
	}
}

func sourceDomain(source string) string {
	return strings.SplitN(source, "/", 2)[0]
}
