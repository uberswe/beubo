package component

import (
	"github.com/markustenghamn/beubo/pkg/structs/page"
	"html/template"
)

type HiddenField struct {
	Theme      string
	Template   string
	Identifier string
	Name       string
	Value      string
	T          *template.Template
}

func (hf HiddenField) GetSection() string {
	return ""
}

func (hf HiddenField) GetTemplateName() string {
	return returnTIfNotEmpty(hf.Template, "component.hiddenfield")
}

func (hf HiddenField) GetTheme() string {
	return returnTIfNotEmpty(hf.Template, "default")
}

func (hf HiddenField) GetTemplate() *template.Template {
	return hf.T
}

func (hf HiddenField) Render() string {
	return page.RenderCompnent(hf)
}
