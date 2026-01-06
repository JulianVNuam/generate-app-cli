package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nuam/scaffold/internal/catalog"
	"github.com/nuam/scaffold/internal/engine"
	"github.com/nuam/scaffold/internal/wizard"
)

func main() {
	// Assumes you run this binary from repo root:
	// repo-root/
	//   cli/
	//   repo-templates/
	repoRoot, err := os.Getwd()
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	templatesRoot := filepath.Join(repoRoot, "repo-templates")
	catalogPath := filepath.Join(templatesRoot, "manifests", "catalog.json")

	cat, err := catalog.LoadCatalog(catalogPath, templatesRoot)
	if err != nil {
		fmt.Println("error loading catalog:", err)
		os.Exit(1)
	}

	plan, err := wizard.Run(cat)
	if err != nil {
		fmt.Println("wizard error:", err)
		os.Exit(1)
	}

	if err := engine.ValidatePlan(cat, plan); err != nil {
		fmt.Println("validation error:", err)
		os.Exit(1)
	}

	if err := engine.ScaffoldProject(cat, templatesRoot, plan); err != nil {
		fmt.Println("scaffold error:", err)
		os.Exit(1)
	}

	fmt.Println("Project generated successfully:", plan.OutputDir)
}
