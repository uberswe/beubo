package menu

import (
	"bytes"
	"fmt"
	"github.com/markustenghamn/beubo/pkg/structs/page"
	"html/template"
	"log"
)

type DefaultMenu struct {
	Items      []page.MenuItem
	Identifier string
	Template   string
	Theme      string
	T          *template.Template
}

func (m DefaultMenu) GetIdentifier() string {
	return m.Identifier
}

func (m DefaultMenu) GetItems() []page.MenuItem {
	return m.Items
}

func (m DefaultMenu) SetItems(items []page.MenuItem) {
	m.Items = items
}

func (m DefaultMenu) Render() string {
	tmpl := "menu.default"
	if m.Template != "" {
		tmpl = m.Template
	}
	theme := "default"
	if m.Theme != "" {
		theme = m.Theme
	}
	path := fmt.Sprintf("%s.%s", theme, tmpl)
	var foundTemplate *template.Template
	if foundTemplate = m.T.Lookup(path); foundTemplate == nil {
		log.Printf("Menu file not found %s\n", path)
		return ""
	}
	buf := &bytes.Buffer{}
	err := foundTemplate.Execute(buf, m)
	if err != nil {
		log.Printf("Component file error executing template %s\n", path)
		return ""
	}
	return buf.String()
}
