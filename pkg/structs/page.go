package structs

import (
	"fmt"
	"github.com/uberswe/beubo/pkg/structs/page"
	"gorm.io/gorm"
	"html/template"
	"log"
)

// Page represents the content of a page, I wanted to go with the concept of having everything be a post even if it's a page, contact form or product
type Page struct {
	gorm.Model
	Title       string `gorm:"size:255"`
	Content     string `sql:"type:text"`
	Description string `sql:"type:text"`
	Excerpt     string `sql:"type:text"`
	Slug        string `gorm:"size:255;unique_index:idx_slug_site_id"`
	Template    string `gorm:"size:255"`
	Site        Site
	SiteID      int   `gorm:"unique_index:idx_slug_site_id"`
	Tags        []Tag `gorm:"many2many:page_tags;"`
}

// PageData is a general structure that holds all data that can be displayed on a page
// using go html templates
type PageData struct {
	Theme       string
	Template    string
	Templates   map[string]string
	Themes      map[string]string
	Title       string
	WebsiteName string
	URL         string
	Error       string
	Warning     string
	Message     string
	Year        string
	Stylesheets []string
	Scripts     []string
	Favicon     string
	Extra       interface{}
	Components  []page.Component
	Menus       []MenuSection
}

// Tag of a post can be used for post categories or things like meta tag keywords for example
type Tag struct {
	gorm.Model
	Value string `gorm:"unique;not null"`
}

// JSONTag is used for responses where tags are shown
// TODO is this redundant if we can use Tag?
type JSONTag struct {
	Value string `json:"value"`
}

// Comment can be related to a post created by a user
type Comment struct {
	gorm.Model
	User    User
	UserID  int
	Email   string
	Website string
	Text    string
	Page    Page
	PageID  int
}

// CreatePage is a method which creates a page using gorm
func CreatePage(db *gorm.DB, title string, slug string, tags []Tag, template string, content string, siteID int) bool {
	pageData := Page{
		Title:    title,
		Content:  content,
		Slug:     slug,
		Template: template,
		SiteID:   siteID,
		Tags:     tags,
	}

	if err := db.Create(&pageData).Error; err != nil {
		fmt.Println("Could not create pageData")
		return false
	}
	return true
}

// FetchPage gets a page based on the provided id from the database
func FetchPage(db *gorm.DB, id int) Page {
	pageData := Page{}

	db.Preload("Tags").First(&pageData, id)

	return pageData
}

// FetchPageBySiteIDAndSlug gets a page based on the site id and slug from the database
func FetchPageBySiteIDAndSlug(db *gorm.DB, SiteID int, slug string) Page {
	pageData := Page{}
	db.Where("slug = ? AND site_id = ?", slug, SiteID).First(&pageData)
	return pageData
}

// UpdatePage is a method which updates a page in the database with relevant data
func UpdatePage(db *gorm.DB, id int, title string, slug string, tags []Tag, template string, content string, siteID int) bool {
	pageData := FetchPage(db, id)

	db.Model(&pageData).Association("Tags").Clear()

	pageData.Title = title
	pageData.Slug = slug
	pageData.Content = content
	pageData.Template = template
	pageData.SiteID = siteID
	pageData.Tags = tags

	if err := db.Save(&pageData).Error; err != nil {
		fmt.Println("Could not create site")
		return false
	}
	return true
}

// DeletePage deletes a page with the provided id from the database
func DeletePage(db *gorm.DB, id int) Page {
	pageData := FetchPage(db, id)
	db.Delete(&pageData)
	return pageData
}

// Content renders components for a page for the specified section
func (pd PageData) Content(section string) template.HTML {
	result := ""
	for _, component := range pd.Components {
		if component.GetSection() == section {
			result += component.Render()
		}
	}
	return template.HTML(result)
}

// Menu renders a menu for the provided section
func (pd PageData) Menu(section string) template.HTML {
	result := ""
	log.Println(section)
	for _, menu := range pd.Menus {
		log.Printf("%s == %s", menu.GetIdentifier(), section)
		if menu.GetIdentifier() == section {
			return template.HTML(menu.Render())
		}
	}
	return template.HTML(result)
}
