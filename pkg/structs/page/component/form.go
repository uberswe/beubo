package component

import (
	"github.com/markustenghamn/beubo/pkg/structs/page"
	"html/template"
)

// Form is a beubo component that can be rendered using HTML templates
type Form struct {
	Section  string
	Fields   []page.Component
	Theme    string
	Template string
	T        *template.Template
	Method   string
	Action   string
}

// GetSection is a getter for the Section property
func (f Form) GetSection() string {
	return f.Section
}

// GetTemplateName is a getter for the Template Property
func (f Form) GetTemplateName() string {
	return returnTIfNotEmpty(f.Template, "component.form")
}

// GetTheme is a getter for the Theme Property
func (f Form) GetTheme() string {
	return returnTIfNotEmpty(f.Theme, "default")
}

// GetTemplate is a getter for the T Property
func (f Form) GetTemplate() *template.Template {
	return f.T
}

// Render calls RenderComponent to turn a Component into a html string for browser output
func (f Form) Render() string {
	return page.RenderComponent(f)
}
