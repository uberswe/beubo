package models

import "github.com/jinzhu/gorm"

type Theme struct {
	gorm.Model
	Title string `gorm:"size:255;unique_index"`
	Slug  string `gorm:"size:255;unique_index"`
}
