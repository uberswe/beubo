package component

import (
	"github.com/uberswe/beubo/pkg/structs/page"
	"html/template"
)

// CheckBoxField is a beubo component that can be rendered using HTML templates
type CheckBoxField struct {
	Theme      string
	Template   string
	Content    string
	Class      string
	Name       string
	Identifier string
	Value      string
	Checked    bool
	T          *template.Template
}

// GetT gets the template.Template for the component
func (cb CheckBoxField) GetT() *template.Template {
	return cb.T
}

// SetT sets the template.Template for the component
func (cb CheckBoxField) SetT(t *template.Template) {
	cb.T = t
}

// GetSection is a getter for the Section property
func (cb CheckBoxField) GetSection() string {
	return ""
}

// GetTemplateName is a getter for the Template Property
func (cb CheckBoxField) GetTemplateName() string {
	return returnTIfNotEmpty(cb.Template, "component.checkboxfield")
}

// GetTheme is a getter for the Theme Property
func (cb CheckBoxField) GetTheme() string {
	return returnTIfNotEmpty(cb.Theme, "default")
}

// GetTemplate is a getter for the T Property
func (cb CheckBoxField) GetTemplate() *template.Template {
	return cb.T
}

// Render calls RenderComponent to turn a Component into a html string for browser output
func (cb CheckBoxField) Render(t *template.Template) string {
	return page.RenderComponent(cb, t)
}
