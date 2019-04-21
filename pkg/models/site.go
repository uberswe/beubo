package models

import (
	"github.com/jinzhu/gorm"
)

type Site struct {
	gorm.Model
	Title     string `gorm:"size:255"`
	Domain    string `gorm:"size:255;unique_index"`
	HandleSsl bool
	Theme     Theme
	ThemeID   int
}
