package structs

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/markustenghamn/beubo/pkg/structs/page"
	"html/template"
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
	Menus       []page.Menu
}

type Tag struct {
	gorm.Model
	Value string `gorm:"unique;not null"`
}

type JsonTag struct {
	Value string `json:"value"`
}

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

// CreateUser is a method which creates a user using gorm
func CreatePage(db *gorm.DB, title string, slug string, tags []Tag, template string, content string, siteID int) bool {
	pageData := Page{
		Title:    title,
		Content:  content,
		Slug:     slug,
		Template: template,
		SiteID:   siteID,
		Tags:     tags,
	}

	if db.NewRecord(pageData) { // => returns `true` as primary key is blank
		if err := db.Create(&pageData).Error; err != nil {
			fmt.Println("Could not create pageData")
			return false
		}
		return true
	}
	return false
}

func FetchPage(db *gorm.DB, id int) Page {
	pageData := Page{}

	db.Preload("Tags").First(&pageData, id)

	return pageData
}

func FetchPageBySiteIDAndSlug(db *gorm.DB, SiteID int, slug string) Page {
	pageData := Page{}

	db.Where("slug = ? AND site_id = ?", slug, SiteID).First(&pageData)

	return pageData
}

// CreateUser is a method which creates a user using gorm
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

func DeletePage(db *gorm.DB, id int) Page {
	pageData := FetchPage(db, id)

	db.Delete(pageData)

	return pageData
}

func (pd PageData) Content(section string) template.HTML {
	result := ""
	for _, component := range pd.Components {
		if component.GetSection() == section {
			result += component.Render()
		}
	}
	return template.HTML(result)
}

func (pd PageData) Menu(section string) template.HTML {
	result := ""
	for _, menu := range pd.Menus {
		if menu.GetIdentifier() == section {
			return template.HTML(menu.Render())
		}
	}
	return template.HTML(result)
}
