package service

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

func ReplacePlaceholders(input string, replacements map[string]string) string {
	re := regexp.MustCompile(`\{\{([^}]+)\}\}`)

	return re.ReplaceAllStringFunc(input, func(match string) string {
		key := strings.TrimSpace(match[2 : len(match)-2])
		if val, ok := replacements[key]; ok {
			return val
		}
		return match
	})
}

func ReplacePathPlaceholders(input string, service string) (*string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	placeholder := map[string]string{
		"home":    home,
		"base":    path.Join(home, "Projects"),
		"service": path.Join(home, "Projects", service),
	}

	result := ReplacePlaceholders(input, placeholder)
	return &result, nil
}

func ReplaceAll(input string, service string, secrets map[string]string) (*string, error) {
	tmp, err := ReplacePathPlaceholders(input, service)
	if err != nil {
		return nil, err
	}

	result := ReplacePlaceholders(*tmp, secrets)
	return &result, nil
}
