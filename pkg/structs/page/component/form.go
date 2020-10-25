package component

import (
	"bytes"
	"fmt"
	"github.com/markustenghamn/beubo/pkg/structs/page"
	"html/template"
	"log"
)

type Form struct {
	Section  string
	Fields   []page.Field
	Theme    string
	Template string
	T        *template.Template
	Method   string
	Action   string
}

func (f Form) GetSection() string {
	return f.Section
}

func (f Form) Render() string {
	tmpl := "component.form"
	if f.Template != "" {
		tmpl = f.Template
	}
	theme := "default"
	if f.Theme != "" {
		theme = f.Theme
	}
	path := fmt.Sprintf("%s.%s", theme, tmpl)
	var foundTemplate *template.Template
	if foundTemplate = f.T.Lookup(path); foundTemplate == nil {
		log.Printf("Component file not found %s\n", path)
		return ""
	}
	buf := &bytes.Buffer{}
	err := foundTemplate.Execute(buf, f)
	if err != nil {
		log.Printf("Component file error executing template %s\n", path)
		return ""
	}
	return buf.String()
}
