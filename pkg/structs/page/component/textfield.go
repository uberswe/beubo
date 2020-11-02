package component

import (
	"github.com/markustenghamn/beubo/pkg/structs/page"
	"html/template"
)

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

func (t TextField) GetSection() string {
	return ""
}

func (t TextField) GetTemplateName() string {
	return returnTIfNotEmpty(t.Template, "component.textfield")
}

func (t TextField) GetTheme() string {
	return returnTIfNotEmpty(t.Template, "default")
}

func (t TextField) GetTemplate() *template.Template {
	return t.T
}

func (t TextField) Render() string {
	return page.RenderCompnent(t)
}
