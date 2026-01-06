package engine

import (
	"os"
	"strings"
)

func GenerateInstructionsBase(basePath, outPath, templateVersion string, modules []string) error {
	b, err := os.ReadFile(basePath)
	if err != nil {
		return err
	}
	content := string(b)
	content = strings.ReplaceAll(content, "__TEMPLATE_VERSION__", templateVersion)

	if len(modules) == 0 {
		content = strings.ReplaceAll(content, "__MODULES_LIST__", "- (none)\n")
	} else {
		var sb strings.Builder
		for _, m := range modules {
			sb.WriteString("- ")
			sb.WriteString(m)
			sb.WriteString("\n")
		}
		content = strings.ReplaceAll(content, "__MODULES_LIST__", sb.String())
	}

	return os.WriteFile(outPath, []byte(content), 0644)
}

func MarkerInsert(filePath string, markerSection string, fragment string) error {
	// marker format expected: <!-- MODULE:<Section> -->
	marker := "<!-- MODULE:" + markerSection + " -->"

	b, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	s := string(b)

	if !strings.Contains(s, marker) {
		// If marker not found, append at end
		s = s + "\n\n" + fragment + "\n"
		return os.WriteFile(filePath, []byte(s), 0644)
	}

	parts := strings.Split(s, marker)
	// Insert below marker
	out := parts[0] + marker + "\n\n" + fragment + "\n" + strings.Join(parts[1:], marker)
	return os.WriteFile(filePath, []byte(out), 0644)
}

func ParseFragmentTarget(fragment string) (target string, body string) {
	// fragment header: <!-- TARGET:Styling -->
	lines := strings.Split(fragment, "\n")
	if len(lines) == 0 {
		return "", fragment
	}
	first := strings.TrimSpace(lines[0])
	if strings.HasPrefix(first, "<!-- TARGET:") && strings.HasSuffix(first, " -->") {
		target = strings.TrimSuffix(strings.TrimPrefix(first, "<!-- TARGET:"), " -->")
		body = strings.Join(lines[1:], "\n")
		body = strings.TrimSpace(body)
		return strings.TrimSpace(target), body
	}
	return "", strings.TrimSpace(fragment)
}

func ReplaceTokenInFile(path, token, value string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	s := strings.ReplaceAll(string(b), token, value)
	return os.WriteFile(path, []byte(s), 0644)
}
