package component

import (
	"github.com/markustenghamn/beubo/pkg/structs/page"
	"html/template"
)

type TextAreaField struct {
	Content    string
	Theme      string
	Template   string
	Class      string
	Identifier string
	Label      string
	Name       string
	Rows       int
	Cols       int
	T          *template.Template
}

func (t TextAreaField) GetSection() string {
	return ""
}

func (t TextAreaField) GetTemplateName() string {
	return returnTIfNotEmpty(t.Template, "component.textareafield")
}

func (t TextAreaField) GetTheme() string {
	return returnTIfNotEmpty(t.Template, "default")
}

func (t TextAreaField) GetTemplate() *template.Template {
	return t.T
}

func (t TextAreaField) Render() string {
	return page.RenderCompnent(t)
}
