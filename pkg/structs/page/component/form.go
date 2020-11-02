package component

import (
	"github.com/markustenghamn/beubo/pkg/structs/page"
	"html/template"
)

type Form struct {
	Section  string
	Fields   []page.Component
	Theme    string
	Template string
	T        *template.Template
	Method   string
	Action   string
}

func (f Form) GetSection() string {
	return f.Section
}

func (f Form) GetTemplateName() string {
	return returnTIfNotEmpty(f.Template, "component.form")
}

func (f Form) GetTheme() string {
	return returnTIfNotEmpty(f.Template, "default")
}

func (f Form) GetTemplate() *template.Template {
	return f.T
}

func (f Form) Render() string {
	return page.RenderCompnent(f)
}
