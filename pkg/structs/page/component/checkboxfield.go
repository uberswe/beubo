package component

import (
	"github.com/markustenghamn/beubo/pkg/structs/page"
	"html/template"
)

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

func (cb CheckBoxField) GetSection() string {
	return ""
}

func (cb CheckBoxField) GetTemplateName() string {
	return returnTIfNotEmpty(cb.Template, "component.checkboxfield")
}

func (cb CheckBoxField) GetTheme() string {
	return returnTIfNotEmpty(cb.Template, "default")
}

func (cb CheckBoxField) GetTemplate() *template.Template {
	return cb.T
}

func (cb CheckBoxField) Render() string {
	return page.RenderCompnent(cb)
}
