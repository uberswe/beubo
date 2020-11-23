package component

import (
	"github.com/markustenghamn/beubo/pkg/structs/page"
	"html/template"
)

// SelectField is a beubo component that can be rendered using HTML templates
type SelectField struct {
	Theme      string
	Template   string
	Identifier string
	Class      string
	Name       string
	Options    []SelectFieldOption
	T          *template.Template
}

// SelectFieldOption is part of the SelectField values, there should be one or more of these
type SelectFieldOption struct {
	Value   string
	Content string
}

// GetSection is a getter for the Section property
func (sf SelectField) GetSection() string {
	return ""
}

// GetTemplateName is a getter for the Template Property
func (sf SelectField) GetTemplateName() string {
	return returnTIfNotEmpty(sf.Template, "component.selectfield")
}

// GetTheme is a getter for the Theme Property
func (sf SelectField) GetTheme() string {
	return returnTIfNotEmpty(sf.Theme, "default")
}

// GetTemplate is a getter for the T Property
func (sf SelectField) GetTemplate() *template.Template {
	return sf.T
}

// Render calls RenderComponent to turn a Component into a html string for browser output
func (sf SelectField) Render() string {
	return page.RenderComponent(sf)
}
