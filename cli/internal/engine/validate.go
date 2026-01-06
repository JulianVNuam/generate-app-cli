package engine

import (
	"fmt"

	"github.com/nuam/scaffold/internal/catalog"
	"github.com/nuam/scaffold/internal/wizard"
)

func ValidatePlan(cat *catalog.Catalog, plan wizard.Plan) error {
	if plan.Stack == "" {
		return fmt.Errorf("stack is required")
	}
	if !cat.CoreExists(plan.Stack) {
		return fmt.Errorf("unknown core stack: %s", plan.Stack)
	}
	// Validate module applicability and basic conflicts by reading manifests
	selected := map[string]bool{}
	for _, mid := range plan.Modules {
		if !cat.ModuleAppliesTo(mid, plan.Stack) {
			return fmt.Errorf("module %s does not apply to stack %s", mid, plan.Stack)
		}
		if selected[mid] {
			return fmt.Errorf("module duplicated: %s", mid)
		}
		selected[mid] = true
	}

	// Conflicts/requires validation (manifest-level)
	for _, mid := range plan.Modules {
		mm, _, err := cat.LoadModuleManifest(plan.Stack, mid)
		if err != nil {
			return fmt.Errorf("cannot load module manifest %s: %w", mid, err)
		}
		for _, c := range mm.Conflicts {
			if selected[c] {
				return fmt.Errorf("module conflict: %s conflicts with %s", mid, c)
			}
		}
		for _, r := range mm.Requires {
			if !selected[r] {
				return fmt.Errorf("module require: %s requires %s", mid, r)
			}
		}
	}
	return nil
}
