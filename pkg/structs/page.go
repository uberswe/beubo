package structs

import (
	"github.com/jinzhu/gorm"
)

// Page represents the content of a page, I wanted to go with the concept of having everything be a post even if it's a page, contact form or product
type Page struct {
	gorm.Model
	Title       string `gorm:"size:255;unique_index"`
	Content     string `sql:"type:text"`
	Description string `sql:"type:text"`
	Excerpt     string `sql:"type:text"`
	Slug        string `gorm:"size:255;unique_index"`
	Template    string `gorm:"size:255"`
	Site        Site
	SiteID      int
}
