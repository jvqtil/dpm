package main

import (
	"fmt"
	"strings"
)

func normalizeSource(source string) string {
	s := strings.TrimSpace(source)
	s = strings.TrimPrefix(s, "https://")
	s = strings.TrimPrefix(s, "http://")
	s = strings.TrimSuffix(s, "/")
	s = strings.TrimSuffix(s, ".git")

	parts := strings.Split(s, "/")

	switch len(parts) {
	case 2:
		return fmt.Sprintf("%s/%s", cfg.AssumeSourceType, s)
	default:
		return s
	}
}

func getSourceDomain(source string) string {
	s := strings.TrimSpace(source)
	s = strings.TrimPrefix(s, "https://")
	s = strings.TrimPrefix(s, "http://")

	return strings.SplitN(s, "/", 2)[0]
}
