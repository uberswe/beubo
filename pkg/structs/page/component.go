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
	Render() string
	GetTemplateName() string
	GetTheme() string
	GetTemplate() *template.Template
}

type ComponentDefinition struct {
	Struct     Component
	Parameters map[string]ComponentParameterDefinition
}

type ComponentParameterDefinition struct {
	StaticValue string
	StructField string
}

func RenderCompnent(c Component) string {
	path := fmt.Sprintf("%s.%s", c.GetTheme(), c.GetTemplateName())
	var foundTemplate *template.Template
	if c.GetTemplate() == nil {
		return ""
	}
	if foundTemplate = c.GetTemplate().Lookup(path); foundTemplate == nil {
		log.Printf("Component file not found %s\n", path)
		return ""
	}
	buf := &bytes.Buffer{}
	err := foundTemplate.Execute(buf, c)
	if err != nil {
		log.Printf("Component file error executing template %s\n", path)
		return ""
	}
	return buf.String()
}
