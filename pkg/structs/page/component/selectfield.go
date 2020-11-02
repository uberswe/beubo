package component

import (
	"github.com/markustenghamn/beubo/pkg/structs/page"
	"html/template"
)

type SelectField struct {
	Theme      string
	Template   string
	Identifier string
	Class      string
	Name       string
	Options    []SelectFieldOption
	T          *template.Template
}

type SelectFieldOption struct {
	Value   string
	Content string
}

func (sf SelectField) GetSection() string {
	return ""
}

func (sf SelectField) GetTemplateName() string {
	return returnTIfNotEmpty(sf.Template, "component.selectfield")
}

func (sf SelectField) GetTheme() string {
	return returnTIfNotEmpty(sf.Template, "default")
}

func (sf SelectField) GetTemplate() *template.Template {
	return sf.T
}

func (sf SelectField) Render() string {
	return page.RenderCompnent(sf)
}
