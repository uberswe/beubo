package page

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
)

// Menu is a component but requires a slice of menu items
type Menu interface {
	GetIdentifier() string
	GetItems() []MenuItem
	SetItems([]MenuItem)
	Render() string
}

// MenuItem is part of a Menu and usually represents a clickable link
type MenuItem struct {
	Text     string
	URI      string
	Template string
	Theme    string
	Items    []MenuItem // A menu item can contain submenus
	T        *template.Template
}

// SubMenu renders submenu items recursively in templates
func (m MenuItem) SubMenu() template.HTML {
	if len(m.Items) > 0 {
		tmpl := "menu.sub"
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
		return template.HTML(buf.String())
	}
	return ""
}
