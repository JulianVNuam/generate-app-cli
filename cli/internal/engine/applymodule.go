package engine

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nuam/scaffold/internal/catalog"
)

func ApplyModule(cat *catalog.Catalog, templatesRoot, stack, moduleID, projectDir string) error {
	mm, manifestPath, err := cat.LoadModuleManifest(stack, moduleID)
	if err != nil {
		return err
	}

	// Module root directory is the directory containing manifest.json
	moduleRoot := filepath.Dir(manifestPath)

	// 1) Copy files
	for _, op := range mm.Files {
		src := filepath.Join(moduleRoot, op.From)
		dst := filepath.Join(projectDir, op.To)
		if op.Mode != "merge" {
			return fmt.Errorf("unsupported file mode: %s (module %s)", op.Mode, moduleID)
		}
		if _, err := os.Stat(src); os.IsNotExist(err) {
			// no files folder is allowed
			continue
		}
		if err := CopyDirMerge(src, dst); err != nil {
			return err
		}
	}

	// 2) Merge package.json patches if present in module files (convention)
	// Convention: if module ships "package.patch.json" at moduleRoot/files/package.patch.json
	patchPackage := filepath.Join(moduleRoot, "files", "package.patch.json")
	targetPackage := filepath.Join(projectDir, "package.json")
	if _, err := os.Stat(patchPackage); err == nil {
		if err := MergePackageJSON(targetPackage, patchPackage); err != nil {
			return err
		}
	}

	// 3) Apply doc patches (markerInsert into instructions.md)
	for _, p := range mm.Patches {
		if p.Type != "markerInsert" {
			return fmt.Errorf("unsupported patch type: %s", p.Type)
		}
		if p.File != "instructions.md" {
			return fmt.Errorf("only instructions.md patching supported in MVP, got: %s", p.File)
		}

		fragPath := filepath.Join(moduleRoot, p.ContentFile)
		b, err := os.ReadFile(fragPath)
		if err != nil {
			return err
		}
		target, body := ParseFragmentTarget(string(b))
		if target == "" {
			// fallback to manifest marker if fragment doesn't declare
			target = p.Marker
		}

		instructionsPath := filepath.Join(projectDir, "instructions.md")
		if err := MarkerInsert(instructionsPath, target, body); err != nil {
			return err
		}
	}

	// 4) env vars are handled at project-level aggregation (EnsureEnvExample)
	return nil
}
