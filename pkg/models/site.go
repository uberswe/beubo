package models

import (
	"github.com/jinzhu/gorm"
)

// Site represents one website, the idea is that Beubo handles many websites at the same time, you could then have 100s of sites all on the same platform
type Site struct {
	gorm.Model
	Title     string `gorm:"size:255"`
	Domain    string `gorm:"size:255;unique_index"`
	HandleSsl bool
	Theme     Theme
	ThemeID   int
}
