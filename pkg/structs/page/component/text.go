package component

import (
	"bytes"
	"fmt"
	beuboTemplate "github.com/markustenghamn/beubo/pkg/template"
	"html/template"
	"log"
)

type Text struct {
	Content  template.HTML
	Theme    string
	Template string
	Class    string
	Renderer beuboTemplate.BeuboTemplateRenderer
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
	if foundTemplate = t.Renderer.T.Lookup(path); foundTemplate == nil {
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
