package models

import "github.com/jinzhu/gorm"

type Config struct {
	gorm.Model
	Key   string `gorm:"size:255;unique_index"`
	Value string `sql:"type:text"`
}
