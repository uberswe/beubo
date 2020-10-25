package component

import (
	"bytes"
	"fmt"
	"github.com/markustenghamn/beubo/pkg/structs/page"
	"html/template"
	"log"
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
	Field page.Field
	Value string
}

func (t Table) Render() string {
	tmpl := "component.table"
	if t.Template != "" {
		tmpl = t.Template
	}
	theme := "default"
	if t.Theme != "" {
		theme = t.Theme
	}
	path := fmt.Sprintf("%s.%s", theme, tmpl)
	var foundTemplate *template.Template
	if foundTemplate = t.T.Lookup(path); foundTemplate == nil {
		log.Printf("Component file not found %s\n", path)
		return ""
	}
	buf := &bytes.Buffer{}
	err := foundTemplate.Execute(buf, t)
	if err != nil {
		log.Printf("Component file error executing template %s\n", path)
		return ""
	}
	return buf.String()
}

func (c Column) RenderField(value string, field page.Field) {

}
