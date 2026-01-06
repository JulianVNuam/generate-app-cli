package catalog

type Catalog struct {
	Cores   []CoreRef   `json:"cores"`
	Modules []ModuleRef `json:"modules"`

	TemplatesRoot string `json:"-"`
}

type CoreRef struct {
	ID   string `json:"id"`   // "next"
	Type string `json:"type"` // "frontend"
	Name string `json:"name"`
}

type ModuleRef struct {
	ID        string   `json:"id"`        // "ui/tailwind"
	Category  string   `json:"category"`  // ui, quality
	AppliesTo []string `json:"appliesTo"` // ["next"]
}

type CoreManifest struct {
	ID          string `json:"id"`          // "core/next"
	Stack       string `json:"stack"`       // "next"
	Version     string `json:"version"`     // "1.0.0"
	TemplateDir string `json:"templateDir"` // "template"
	DocsBase    string `json:"docsBase"`    // "docs/instructions.base.md"
	Env         []string `json:"env"`       // base env vars
}

type ModuleManifest struct {
	ID        string   `json:"id"`        // "ui/tailwind"
	Name      string   `json:"name"`
	AppliesTo []string `json:"appliesTo"`
	Requires  []string `json:"requires,omitempty"`
	Conflicts []string `json:"conflicts,omitempty"`

	Files []FileOp `json:"files,omitempty"`
	Deps  Deps     `json:"deps,omitempty"`

	Patches []Patch `json:"patches,omitempty"`
	Env     []string `json:"env,omitempty"`
}

type FileOp struct {
	From string `json:"from"`
	To   string `json:"to"`
	Mode string `json:"mode"` // "merge" supported
}

type Deps struct {
	NPM *NPMDeps `json:"npm,omitempty"`
}

type NPMDeps struct {
	Dependencies    []string `json:"dependencies,omitempty"`
	DevDependencies []string `json:"devDependencies,omitempty"`
	Scripts         map[string]string `json:"scripts,omitempty"`
}

type Patch struct {
	Type        string `json:"type"`        // "markerInsert"
	File        string `json:"file"`        // "instructions.md"
	Marker      string `json:"marker"`      // "Styling", etc.
	ContentFile string `json:"contentFile"` // "docs/instructions.fragment.md"
}
