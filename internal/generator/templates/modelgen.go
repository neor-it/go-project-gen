// internal/generator/templates/modelgen.go - Templates for model generator
package templates

import (
	_ "embed" // Required for go:embed
	// Keep fmt for Sprintf - Removed fmt import as it's no longer used
	_ "github.com/lib/pq"
)

//go:embed tmpls/modelgen_script.tmpl
var modelGeneratorScriptContent string

// ModelGeneratorFullTemplate returns the complete content of the model generator tool,
// read from the embedded modelgen_script.tmpl file.
func ModelGeneratorFullTemplate() string {
	// The content is embedded directly into the modelGeneratorScriptContent variable by the go:embed directive.
	// We no longer need fmt.Sprintf here.
	return modelGeneratorScriptContent
} // End of ModelGeneratorFullTemplate
