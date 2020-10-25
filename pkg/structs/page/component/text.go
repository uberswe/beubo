package component

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
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

func (t Text) Render() string {
	tmpl := "component.text"
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
