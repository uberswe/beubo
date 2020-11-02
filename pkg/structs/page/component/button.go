package component

import (
	"github.com/markustenghamn/beubo/pkg/structs/page"
	"html/template"
)

type Button struct {
	Theme    string
	Template string
	Class    string
	Content  string
	T        *template.Template
}

func (b Button) GetSection() string {
	return ""
}

func (b Button) GetTemplateName() string {
	return returnTIfNotEmpty(b.Template, "component.button")
}

func (b Button) GetTheme() string {
	return returnTIfNotEmpty(b.Theme, "default")
}

func (b Button) GetTemplate() *template.Template {
	return b.T
}

func (b Button) Render() string {
	return page.RenderCompnent(b)
}
