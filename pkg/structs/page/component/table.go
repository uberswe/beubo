package component

import (
	"github.com/markustenghamn/beubo/pkg/structs/page"
	"html/template"
)

type Table struct {
	Section          string
	Header           []Column
	Rows             []Row
	Theme            string
	Template         string
	PageNumber       int // Current page
	PageDisplayCount int // How many rows per page
	T                *template.Template
}

func (t Table) GetSection() string {
	return t.Section
}

type Row struct {
	Columns []Column
}

type Column struct {
	Name  string
	Field page.Component
	Value string
}

func (t Table) GetTemplateName() string {
	return returnTIfNotEmpty(t.Template, "component.table")
}

func (t Table) GetTheme() string {
	return returnTIfNotEmpty(t.Template, "default")
}

func (t Table) GetTemplate() *template.Template {
	return t.T
}

func (t Table) Render() string {
	return page.RenderCompnent(t)
}

func (t Table) RenderColumn(c Column) string {
	return page.RenderCompnent(c.Field)
}

func (c Column) RenderField(value string, field page.Component) string {
	if field != nil && field.Render() != "" {
		return field.Render()
	}
	return value
}
