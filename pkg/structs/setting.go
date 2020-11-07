package structs

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type Setting struct {
	gorm.Model
	Key   string `gorm:"size:255;unique_index"`
	Value string `gorm:"size:255"`
}

// CreateSite is a method which creates a site using gorm
func CreateSetting(db *gorm.DB, key string, value string) bool {
	setting := Setting{
		Key:   key,
		Value: value,
	}

	if db.NewRecord(setting) { // => returns `true` as primary key is blank
		if err := db.Create(&setting).Error; err != nil {
			fmt.Println("Could not create setting")
			return false
		}
		return true
	}
	return false
}

func FetchSetting(db *gorm.DB, id int) Setting {
	setting := Setting{}

	db.First(&setting, id)

	return setting
}

// UpdateSetting updates a setting key value pair using gorm
func UpdateSetting(db *gorm.DB, id int, key string, value string) bool {
	setting := FetchSetting(db, id)

	setting.Key = key
	setting.Value = value

	if err := db.Save(&setting).Error; err != nil {
		fmt.Println("Could not create setting")
		return false
	}
	return true
}

func DeleteSetting(db *gorm.DB, id int) Setting {
	setting := FetchSetting(db, id)

	db.Delete(setting)

	return setting
}
