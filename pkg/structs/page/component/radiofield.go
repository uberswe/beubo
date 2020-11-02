package component

import (
	"github.com/markustenghamn/beubo/pkg/structs/page"
	"html/template"
)

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

func (rf RadioField) GetSection() string {
	return ""
}

func (rf RadioField) GetTemplateName() string {
	return returnTIfNotEmpty(rf.Template, "component.radiofield")
}

func (rf RadioField) GetTheme() string {
	return returnTIfNotEmpty(rf.Template, "default")
}

func (rf RadioField) GetTemplate() *template.Template {
	return rf.T
}

func (rf RadioField) Render() string {
	return page.RenderCompnent(rf)
}
