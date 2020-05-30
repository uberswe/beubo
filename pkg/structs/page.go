package structs

import (
	"fmt"
	"github.com/jinzhu/gorm"
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
	SiteID      int `gorm:"unique_index:idx_slug_site_id"`
}

// PageData is a general structure that holds all data that can be displayed on a page
// using go html templates
type PageData struct {
	Title       string
	WebsiteName string
	URL         string
	Menu        []MenuItem
	Error       string
	Warning     string
	Message     string
	Year        string
	Extra       interface{}
}

// MenuItem is one item that can be part of a nav in the frontend
// TODO might be too specific consider removing/redoing this
type MenuItem struct {
	Title string
	Path  string
}

// CreateUser is a method which creates a user using gorm
func CreatePage(db *gorm.DB, title string, slug string, content string, siteID int) bool {
	page := Page{
		Title:   title,
		Content: content,
		Slug:    slug,
		SiteID:  siteID,
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

	db.First(&page, id)

	return page
}

// CreateUser is a method which creates a user using gorm
func UpdatePage(db *gorm.DB, id int, title string, slug string, content string, siteID int) bool {
	page := FetchPage(db, id)

	page.Title = title
	page.Slug = slug
	page.Content = content
	page.SiteID = siteID

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
