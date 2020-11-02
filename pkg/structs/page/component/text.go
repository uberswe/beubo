package component

import (
	"github.com/markustenghamn/beubo/pkg/structs/page"
	"html/template"
)

type Text struct {
	Section  string
	Content  template.HTML
	Theme    string
	Template string
	Class    string
	T        *template.Template
}

func (t Text) GetSection() string {
	return t.Section
}

func (t Text) GetTemplateName() string {
	return returnTIfNotEmpty(t.Template, "component.text")
}

func (t Text) GetTheme() string {
	return returnTIfNotEmpty(t.Template, "default")
}

func (t Text) GetTemplate() *template.Template {
	return t.T
}

func (t Text) Render() string {
	return page.RenderCompnent(t)
}
