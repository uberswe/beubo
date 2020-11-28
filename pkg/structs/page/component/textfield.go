package component

import (
	"github.com/uberswe/beubo/pkg/structs/page"
	"html/template"
)

// TextField is a beubo component that can be rendered using HTML templates
type TextField struct {
	Theme       string
	Template    string
	Class       string
	Identifier  string
	Label       string
	Name        string
	Value       string
	Placeholder string
	T           *template.Template
}

// GetSection is a getter for the Section property
func (t TextField) GetSection() string {
	return ""
}

// GetTemplateName is a getter for the Template property
func (t TextField) GetTemplateName() string {
	return returnTIfNotEmpty(t.Template, "component.textfield")
}

// GetTheme is a getter for the Theme property
func (t TextField) GetTheme() string {
	return returnTIfNotEmpty(t.Theme, "default")
}

// GetTemplate is a getter for the T Property
func (t TextField) GetTemplate() *template.Template {
	return t.T
}

// Render calls RenderComponent to turn a Component into a html string for browser output
func (t TextField) Render() string {
	return page.RenderComponent(t)
}
