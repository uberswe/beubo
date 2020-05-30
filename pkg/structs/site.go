package structs

import (
	"fmt"
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

// CreateUser is a method which creates a user using gorm
func CreateSite(db *gorm.DB, title string, domain string, ssl bool) bool {
	site := Site{
		Title:     title,
		Domain:    domain,
		HandleSsl: ssl,
	}

	if db.NewRecord(site) { // => returns `true` as primary key is blank
		if err := db.Create(&site).Error; err != nil {
			fmt.Println("Could not create site")
			return false
		}
		return true
	}
	return false
}

func FetchSite(db *gorm.DB, id int) Site {
	site := Site{}

	db.First(&site, id)

	return site
}

// CreateUser is a method which creates a user using gorm
func UpdateSite(db *gorm.DB, id int, title string, domain string, ssl bool) bool {
	site := FetchSite(db, id)

	site.Title = title
	site.Domain = domain
	site.HandleSsl = ssl

	if err := db.Save(&site).Error; err != nil {
		fmt.Println("Could not create site")
		return false
	}
	return true
}

func DeleteSite(db *gorm.DB, id int) Site {
	site := FetchSite(db, id)

	db.Delete(site)

	return site
}
