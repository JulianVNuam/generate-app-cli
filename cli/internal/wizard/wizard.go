package wizard

import (
	"fmt"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nuam/scaffold/internal/catalog"
)

func Run(cat *catalog.Catalog) (Plan, error) {
	var p Plan

	if err := survey.AskOne(&survey.Input{Message: "Project name:"}, &p.ProjectName, survey.WithValidator(survey.Required)); err != nil {
		return p, err
	}

	if err := survey.AskOne(&survey.Select{
		Message: "Project type:",
		Options: []string{"frontend", "backend"},
	}, &p.Kind, survey.WithValidator(survey.Required)); err != nil {
		return p, err
	}

	// For this first test case: only frontend + Next
	if p.Kind == "frontend" {
		p.Stack = "next"
		fmt.Println("Framework:", "next (App Router)")
	} else {
		// user said: only next for now; keep behavior strict
		return p, fmt.Errorf("for now only frontend/next is enabled in this test case")
	}

	// Step 3: dependencies/features
	var modules []string

	// UI choice: Tailwind? (Yes/No)
	var wantTailwind bool
	if err := survey.AskOne(&survey.Confirm{Message: "Include Tailwind CSS?"}, &wantTailwind); err != nil {
		return p, err
	}
	if wantTailwind {
		modules = append(modules, "ui/tailwind")
	}

	// Quality: ESLint+Prettier? (Yes/No)
	var wantQuality bool
	if err := survey.AskOne(&survey.Confirm{Message: "Include ESLint + Prettier?"}, &wantQuality); err != nil {
		return p, err
	}
	if wantQuality {
		modules = append(modules, "quality/eslint-prettier")
	}

	p.Modules = modules

	// Step 4: configuration/best practices (for MVP we keep it folded into modules)
	// Could expand later with additional prompts.

	// Output directory
	p.OutputDir = filepath.Join(".", p.ProjectName)

	return p, nil
}
