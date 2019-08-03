package models

import "github.com/jinzhu/gorm"

// Theme is the template or html files that a site uses, theme files are found under the web directory
type Theme struct {
	gorm.Model
	Title string `gorm:"size:255;unique_index"`
	Slug  string `gorm:"size:255;unique_index"`
}
