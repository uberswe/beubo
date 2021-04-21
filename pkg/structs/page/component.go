package page

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
)

// Component is anything that can be rendered on a page. A text field is a component but so is the form the text field is a part of.
type Component interface {
	GetSection() string
	// render returns a html template string with the content of the field
	Render(t *template.Template) string
	GetTemplateName() string
	GetTheme() string
	SetT(t *template.Template)
	GetT() *template.Template
}

// RenderComponent takes the provided component and finds the relevant template and renders this into a string
func RenderComponent(c Component, t *template.Template) string {
	path := fmt.Sprintf("%s.%s", c.GetTheme(), c.GetTemplateName())
	if t.Lookup(path) == nil {
		log.Printf("Component file not found %s\n", path)
		return ""
	}
	buf := &bytes.Buffer{}
	err := t.Lookup(path).Execute(buf, c)
	if err != nil {
		log.Printf("Component file error executing template %s\n", path)
		return ""
	}
	return buf.String()
}
