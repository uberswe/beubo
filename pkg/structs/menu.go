package structs

import (
	"bytes"
	"fmt"
	"gorm.io/gorm"
	"html/template"
	"log"
)

type MenuSection struct {
	gorm.Model
	Section  string     `gorm:"size:255;unique_index:idx_section_site_id"`
	Items    []MenuItem `gorm:"foreignKey:SectionID"`
	Site     Site
	SiteID   int `gorm:"unique_index:idx_section_site_id"`
	Template string
	Theme    string             `gorm:"-"`
	T        *template.Template `gorm:"-"`
}

func (m MenuSection) GetIdentifier() string {
	return m.Section
}

func (m MenuSection) GetItems() []MenuItem {
	return m.Items
}

func (m MenuSection) Render() string {
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

// MenuItem is part of a Menu and usually represents a clickable link
type MenuItem struct {
	gorm.Model
	SectionID     *int
	ParentID      *int
	Text          string
	URI           string
	Authenticated bool
	Template      string             `gorm:"-"`
	Theme         string             `gorm:"-"`
	Items         []MenuItem         `gorm:"foreignKey:ParentID"` // A menu item can contain submenus
	Permissions   []MenuPermission   `gorm:"foreignKey:MenuItemID"`
	Settings      []MenuSetting      `gorm:"foreignKey:MenuItemID"`
	T             *template.Template `gorm:"-"`
}

// MenuPermission defines permissions needed to show a menu item
type MenuPermission struct {
	gorm.Model
	Permission string
	// If the permission exists should the menu be shown?
	Show       bool
	MenuItem   MenuItem
	MenuItemID int
}

// MenuSetting defines settings needed to show a menu item
type MenuSetting struct {
	gorm.Model
	Setting     string
	ShouldEqual string
	// If the setting exists should the menu be shown?
	Show       bool
	MenuItem   MenuItem
	MenuItemID int
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

func (m MenuItem) RenderEditSection() template.HTML {
	log.Println("RenderEditSection")
	tmpl := "admin.menu.section"
	if m.Template != "" {
		tmpl = m.Template
	}
	theme := "default"
	if m.Theme != "" {
		theme = m.Theme
	}
	path := fmt.Sprintf("%s.%s", theme, tmpl)

	log.Println(path)
	var foundTemplate *template.Template
	if foundTemplate = m.T.Lookup(path); foundTemplate == nil {
		log.Printf("Menu file not found %s\n", path)
		return ""
	}
	log.Println("RenderEditSection 2")
	buf := &bytes.Buffer{}
	err := foundTemplate.Execute(buf, m)
	if err != nil {
		log.Printf("Component file error executing template %s\n", path)
		return ""
	}
	log.Println("RenderEditSection 3")
	return template.HTML(buf.String())
}

func CreateMenu(db *gorm.DB, section string, template string, siteID int) MenuSection {
	menuSection := MenuSection{
		Section:  section,
		Template: template,
		SiteID:   siteID,
	}

	if err := db.Create(&menuSection).Error; err != nil {
		fmt.Println("Could not create menuSection")
		return menuSection
	}
	return menuSection
}

func DeleteMenu(db *gorm.DB, id int, siteID int) MenuSection {
	setting := FetchMenuWithSiteID(db, id, siteID)
	db.Delete(&setting)
	return setting
}

func FetchMenu(db *gorm.DB, id int) MenuSection {
	menuSection := MenuSection{}
	db.Where("id = ?", id).First(&menuSection)
	return menuSection
}

func FetchMenuWithSiteID(db *gorm.DB, id int, siteID int) MenuSection {
	menuSection := MenuSection{}
	db.Where("id = ?", id).Where("site_id = ?", siteID).First(&menuSection)
	return menuSection
}

func FetchMenusBySiteID(db *gorm.DB, siteID int) []MenuSection {
	var menuSections []MenuSection
	db.Where("site_id = ?", siteID).Find(&menuSections)
	return menuSections
}

func FetchMenuItemsBySectionId(db *gorm.DB, sectionID int) []MenuItem {
	var menuItems []MenuItem
	db.Where("section_id = ?", sectionID).Find(&menuItems)
	return menuItems
}

func FetchMenuItemsByParentId(db *gorm.DB, parentID int) []MenuItem {
	var menuItems []MenuItem
	db.Where("parent_id = ?", parentID).Find(&menuItems)
	return menuItems
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
