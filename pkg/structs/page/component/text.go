package component

import (
	"github.com/uberswe/beubo/pkg/structs/page"
	"html/template"
)

// Text is a beubo component that can be rendered using HTML templates
type Text struct {
	Section  string
	Content  template.HTML
	Theme    string
	Template string
	Class    string
	T        *template.Template
}

// GetSection is a getter for the Section property
func (t Text) GetSection() string {
	return t.Section
}

// GetT gets the template.Template for the component
func (t Text) GetT() *template.Template {
	return t.T
}

// SetT sets the template.Template for the component
func (t Text) SetT(temp *template.Template) {
	t.T = temp
}

// GetTemplateName is a getter for the Template property
func (t Text) GetTemplateName() string {
	return returnTIfNotEmpty(t.Template, "component.text")
}

// GetTheme is a getter for the Theme property
func (t Text) GetTheme() string {
	return returnTIfNotEmpty(t.Theme, "default")
}

// GetTemplate is a getter for the T Property
func (t Text) GetTemplate() *template.Template {
	return t.T
}

// Render calls RenderComponent to turn a Component into a html string for browser output
func (t Text) Render(te *template.Template) string {
	return page.RenderComponent(t, te)
}
