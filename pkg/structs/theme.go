package structs

import "github.com/jinzhu/gorm"

// Theme is the template or html files that a site uses, theme files are found under the web directory
type Theme struct {
	gorm.Model
	Title string `gorm:"size:255;unique_index"`
	Slug  string `gorm:"size:255;unique_index"`
}

func FetchTheme(db *gorm.DB, id int) Theme {
	theme := Theme{}

	db.First(&theme, id)

	return theme
}

func FetchThemeBySlug(db *gorm.DB, slug string) Theme {
	theme := Theme{}

	db.Where("slug = ?", slug).First(&theme)

	return theme
}
