package models

import "github.com/jinzhu/gorm"

// Config is a key value store for various settings and configurations
type Config struct {
	gorm.Model
	Key   string `gorm:"size:255;unique_index"`
	Value string `sql:"type:text"`
}
