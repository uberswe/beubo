package component

import (
	"github.com/uberswe/beubo/pkg/structs/page"
	"html/template"
)

// Table is a beubo component that can be rendered using HTML templates
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

// GetSection is a getter for the Section property
func (t Table) GetSection() string {
	return t.Section
}

// Row represents a html table row which can have columns
type Row struct {
	Columns []Column
}

// Column represents a html column in a table which is part of a row
type Column struct {
	Name  string
	Field page.Component
	Value string
}

// GetT gets the template.Template for the component
func (t Table) GetT() *template.Template {
	return t.T
}

// SetT sets the template.Template for the component
func (t Table) SetT(temp *template.Template) {
	t.T = temp
}

// GetTemplateName is a getter for the Template Property
func (t Table) GetTemplateName() string {
	return returnTIfNotEmpty(t.Template, "component.table")
}

// GetTheme is a getter for the Theme Property
func (t Table) GetTheme() string {
	return returnTIfNotEmpty(t.Theme, "default")
}

// GetTemplate is a getter for the T Property
func (t Table) GetTemplate() *template.Template {
	return t.T
}

// Render calls RenderComponent to turn a Component into a html string for browser output
func (t Table) Render(te *template.Template) string {
	return page.RenderComponent(t, te)
}

// RenderColumn calls RenderComponent to turn a Column into a html string which is added to the Table Render
func (t Table) RenderColumn(c Column, te *template.Template) string {
	return page.RenderComponent(c.Field, te)
}

// RenderField calls Render to turn a Column into a string which is added to the Table Render
func (c Column) RenderField(value string, field page.Component, t *template.Template) template.HTML {
	if field != nil && field.Render(t) != "" {
		return template.HTML(field.Render(t))
	}
	return template.HTML(value)
}
