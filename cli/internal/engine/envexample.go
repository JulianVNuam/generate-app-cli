package engine

import (
	"os"
	"path/filepath"
	"strings"
)

func EnsureEnvExample(projectDir string, vars []string) error {
	envPath := filepath.Join(projectDir, ".env.example")
	existingLines := map[string]bool{}

	if b, err := os.ReadFile(envPath); err == nil {
		for _, line := range strings.Split(string(b), "\n") {
			t := strings.TrimSpace(line)
			if t == "" || strings.HasPrefix(t, "#") {
				continue
			}
			existingLines[t] = true
		}
	}

	// Ensure file exists
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		if err := os.WriteFile(envPath, []byte("# Environment variables\n"), 0644); err != nil {
			return err
		}
	}

	var toAppend []string
	for _, v := range vars {
		if strings.TrimSpace(v) == "" {
			continue
		}
		line := v + "="
		if !existingLines[line] {
			toAppend = append(toAppend, line)
		}
	}

	if len(toAppend) == 0 {
		return nil
	}

	f, err := os.OpenFile(envPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, l := range toAppend {
		if _, err := f.WriteString(l + "\n"); err != nil {
			return err
		}
	}

	return nil
}
