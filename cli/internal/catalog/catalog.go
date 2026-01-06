package catalog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func LoadCatalog(catalogPath string, templatesRoot string) (*Catalog, error) {
	b, err := os.ReadFile(catalogPath)
	if err != nil {
		return nil, err
	}
	var c Catalog
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	c.TemplatesRoot = templatesRoot
	return &c, nil
}

func (c *Catalog) CoreExists(stack string) bool {
	for _, core := range c.Cores {
		if core.ID == stack {
			return true
		}
	}
	return false
}

func (c *Catalog) ModuleAppliesTo(moduleID, stack string) bool {
	for _, m := range c.Modules {
		if m.ID == moduleID {
			for _, s := range m.AppliesTo {
				if s == stack {
					return true
				}
			}
		}
	}
	return false
}

func (c *Catalog) LoadCoreManifest(stack string) (*CoreManifest, string, error) {
	manifestPath := filepath.Join(c.TemplatesRoot, "core", stack, "manifest.json")
	b, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, "", err
	}
	var m CoreManifest
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}
	if m.Stack != stack {
		return nil, "", fmt.Errorf("core manifest stack mismatch: expected %s got %s", stack, m.Stack)
	}
	return &m, manifestPath, nil
}

func (c *Catalog) LoadModuleManifest(stack, moduleID string) (*ModuleManifest, string, error) {
	// moduleID example: "ui/tailwind"
	parts := strings.SplitN(moduleID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, "", fmt.Errorf("invalid module id: %s", moduleID)
	}
	cat := parts[0]
	name := parts[1]

	manifestPath := filepath.Join(c.TemplatesRoot, "modules", cat, name, stack, "manifest.json")
	b, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, "", err
	}
	var mm ModuleManifest
	if err := json.Unmarshal(b, &mm); err != nil {
		return nil, "", err
	}
	return &mm, manifestPath, nil
}
