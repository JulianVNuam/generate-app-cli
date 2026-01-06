package wizard

type Plan struct {
	ProjectName string
	Kind        string   // "frontend"|"backend"
	Stack       string   // "next" (for now)
	Modules     []string // e.g. ["ui/tailwind","quality/eslint-prettier"]
	OutputDir   string   // where to generate
}
