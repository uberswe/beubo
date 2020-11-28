package component

import (
	"github.com/uberswe/beubo/pkg/structs/page"
	"html/template"
)

// HiddenField is a beubo component that can be rendered using HTML templates
type HiddenField struct {
	Theme      string
	Template   string
	Identifier string
	Name       string
	Value      string
	T          *template.Template
}

// GetSection is a getter for the Section property
func (hf HiddenField) GetSection() string {
	return ""
}

// GetTemplateName is a getter for the Template Property
func (hf HiddenField) GetTemplateName() string {
	return returnTIfNotEmpty(hf.Template, "component.hiddenfield")
}

// GetTheme is a getter for the Theme Property
func (hf HiddenField) GetTheme() string {
	return returnTIfNotEmpty(hf.Theme, "default")
}

// GetTemplate is a getter for the T Property
func (hf HiddenField) GetTemplate() *template.Template {
	return hf.T
}

// Render calls RenderComponent to turn a Component into a html string for browser output
func (hf HiddenField) Render() string {
	return page.RenderComponent(hf)
}
