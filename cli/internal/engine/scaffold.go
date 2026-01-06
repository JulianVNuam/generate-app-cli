package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nuam/scaffold/internal/catalog"
	"github.com/nuam/scaffold/internal/wizard"
)

func ScaffoldProject(cat *catalog.Catalog, templatesRoot string, plan wizard.Plan) error {
	// 1) Load core manifest
	core, _, err := cat.LoadCoreManifest(plan.Stack)
	if err != nil {
		return err
	}

	// 2) Ensure output dir does not exist
	if _, err := os.Stat(plan.OutputDir); err == nil {
		return fmt.Errorf("output directory already exists: %s", plan.OutputDir)
	}

	// 3) Copy core template
	coreTemplatePath := filepath.Join(templatesRoot, "core", plan.Stack, core.TemplateDir)
	if err := CopyDirMerge(coreTemplatePath, plan.OutputDir); err != nil {
		return err
	}

	// 4) Generate instructions.md from base
	baseDocPath := filepath.Join(templatesRoot, "core", plan.Stack, core.DocsBase)
	instructionsOut := filepath.Join(plan.OutputDir, "instructions.md")
	if err := GenerateInstructionsBase(baseDocPath, instructionsOut, core.Version, plan.Modules); err != nil {
		return err
	}

	// 5) Apply selected modules
	var allEnv []string
	allEnv = append(allEnv, core.Env...)

	for _, mid := range plan.Modules {
		if err := ApplyModule(cat, templatesRoot, plan.Stack, mid, plan.OutputDir); err != nil {
			return err
		}

		mm, _, err := cat.LoadModuleManifest(plan.Stack, mid)
		if err != nil {
			return err
		}
		allEnv = append(allEnv, mm.Env...)
	}

	// 6) Build .env.example from collected env vars
	allEnv = unique(allEnv)
	if err := EnsureEnvExample(plan.OutputDir, allEnv); err != nil {
		return err
	}

	// 7) Replace project name tokens (optional)
	_ = ReplaceTokenInFile(filepath.Join(plan.OutputDir, "package.json"), "__PROJECT_NAME__", plan.ProjectName)

	return nil
}

func unique(in []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(in))
	for _, v := range in {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}
