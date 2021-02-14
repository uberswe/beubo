package page

import (
	"bytes"
	"fmt"
	"gorm.io/gorm"
	"html/template"
	"log"
)

// Menu is a component but requires a slice of menu items
type Menu interface {
	GetIdentifier() string
	GetItems() []MenuItem
	Render() string
}

type MenuSection struct {
	gorm.Model
	Section string
	Items   []MenuItem `gorm:"foreignKey:SectionID"`
}

// MenuItem is part of a Menu and usually represents a clickable link
type MenuItem struct {
	gorm.Model
	SectionID     uint
	Text          string
	URI           string
	Authenticated bool
	Always        bool
	Template      string     `gorm:"-"`
	Theme         string     `gorm:"-"`
	Items         []MenuItem `gorm:"foreignKey:SectionID"` // A menu item can contain submenus
	// TODO site id can be nullable when this happens it's a global menu
	SiteID uint
	T      *template.Template `gorm:"-"`
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

func CreateMenu(db *gorm.DB, section string) MenuSection {
	menuSection := MenuSection{
		Section: section,
	}

	if err := db.Create(&menuSection).Error; err != nil {
		fmt.Println("Could not create menuSection")
		return menuSection
	}
	return menuSection
}

func DeleteMenu(db *gorm.DB, id int) MenuSection {
	setting := FetchMenu(db, id)
	db.Delete(&setting)
	return setting
}

func FetchMenu(db *gorm.DB, id int) MenuSection {
	menuSection := MenuSection{}
	db.Where("id = ?", id).First(&menuSection)
	return menuSection
}

func UpdateMenu(db *gorm.DB, id int, section string) bool {
	menuSection := FetchMenu(db, id)

	menuSection.Section = section

	if err := db.Save(&menuSection).Error; err != nil {
		fmt.Println("Could not create menuSetting")
		return false
	}
	return true
}
