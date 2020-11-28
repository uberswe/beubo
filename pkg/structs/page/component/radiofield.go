package component

import (
	"github.com/uberswe/beubo/pkg/structs/page"
	"html/template"
)

// RadioField is a beubo component that can be rendered using HTML templates
type RadioField struct {
	Theme      string
	Template   string
	Class      string
	Identifier string
	Name       string
	Value      string
	Content    string
	Checked    bool
	T          *template.Template
}

// GetSection is a getter for the Section property
func (rf RadioField) GetSection() string {
	return ""
}

// GetTemplateName is a getter for the Template Property
func (rf RadioField) GetTemplateName() string {
	return returnTIfNotEmpty(rf.Template, "component.radiofield")
}

// GetTheme is a getter for the Theme Property
func (rf RadioField) GetTheme() string {
	return returnTIfNotEmpty(rf.Theme, "default")
}

// GetTemplate is a getter for the T Property
func (rf RadioField) GetTemplate() *template.Template {
	return rf.T
}

// Render calls RenderComponent to turn a Component into a html string for browser output
func (rf RadioField) Render() string {
	return page.RenderComponent(rf)
}
