package structs

import "github.com/jinzhu/gorm"

type Setting struct {
	gorm.Model
	Key   string `gorm:"size:255;unique_index"`
	Value string `gorm:"size:255"`
}
