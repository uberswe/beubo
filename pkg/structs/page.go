package structs

import (
	"fmt"
	"github.com/jinzhu/gorm"
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
	Content     template.HTML
	WebsiteName string
	URL         string
	Menu        []MenuItem
	Error       string
	Warning     string
	Message     string
	Year        string
	Stylesheets []string
	Scripts     []string
	Favicon     string
	Extra       interface{}
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

// MenuItem is one item that can be part of a nav in the frontend
// TODO might be too specific consider removing/redoing this
type MenuItem struct {
	Title string
	Path  string
}

// CreateUser is a method which creates a user using gorm
func CreatePage(db *gorm.DB, title string, slug string, tags []Tag, template string, content string, siteID int) bool {
	page := Page{
		Title:    title,
		Content:  content,
		Slug:     slug,
		Template: template,
		SiteID:   siteID,
		Tags:     tags,
	}

	if db.NewRecord(page) { // => returns `true` as primary key is blank
		if err := db.Create(&page).Error; err != nil {
			fmt.Println("Could not create page")
			return false
		}
		return true
	}
	return false
}

func FetchPage(db *gorm.DB, id int) Page {
	page := Page{}

	db.Preload("Tags").First(&page, id)

	return page
}

func FetchPageBySiteIDAndSlug(db *gorm.DB, SiteID int, slug string) Page {
	page := Page{}

	db.Where("slug = ? AND site_id = ?", slug, SiteID).First(&page)

	return page
}

// CreateUser is a method which creates a user using gorm
func UpdatePage(db *gorm.DB, id int, title string, slug string, tags []Tag, template string, content string, siteID int) bool {
	page := FetchPage(db, id)

	db.Model(&page).Association("Tags").Clear()

	page.Title = title
	page.Slug = slug
	page.Content = content
	page.Template = template
	page.SiteID = siteID
	page.Tags = tags

	if err := db.Save(&page).Error; err != nil {
		fmt.Println("Could not create site")
		return false
	}
	return true
}

func DeletePage(db *gorm.DB, id int) Page {
	page := FetchPage(db, id)

	db.Delete(page)

	return page
}
