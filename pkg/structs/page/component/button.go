package component

import (
	"github.com/uberswe/beubo/pkg/structs/page"
	"html/template"
)

// Button is a beubo component that can be rendered using HTML templates
type Button struct {
	Section  string
	Theme    string
	Template string
	Class    string
	Content  string
	Link     template.URL
	T        *template.Template
}

// GetT gets the template.Template for the component
func (b Button) GetT() *template.Template {
	return b.T
}

// SetT sets the template.Template for the component
func (b Button) SetT(t *template.Template) {
	b.T = t
}

// GetSection is a getter for the Section property
func (b Button) GetSection() string {
	return b.Section
}

// GetTemplateName is a getter for the Template Property
func (b Button) GetTemplateName() string {
	return returnTIfNotEmpty(b.Template, "component.button")
}

// GetTheme is a getter for the Theme Property
func (b Button) GetTheme() string {
	return returnTIfNotEmpty(b.Theme, "default")
}

// GetTemplate is a getter for the T Property
func (b Button) GetTemplate() *template.Template {
	return b.T
}

// Render calls RenderComponent to turn a Component into a html string for browser output
func (b Button) Render(t *template.Template) string {
	return page.RenderComponent(b, t)
}
