package structs

import (
	"fmt"
	"gorm.io/gorm"
)

// Setting represents a key value setting for Beubo usually used for global config values
type Setting struct {
	gorm.Model
	Key   string `gorm:"size:255"`
	Value string `gorm:"size:255"`
}

// CreateSetting is a method which creates a setting using gorm
func CreateSetting(db *gorm.DB, key string, value string) bool {
	setting := Setting{
		Key:   key,
		Value: value,
	}

	if err := db.Create(&setting).Error; err != nil {
		fmt.Println("Could not create setting")
		return false
	}
	return true
}

// FetchSetting gets a setting from the database via the provided id
func FetchSetting(db *gorm.DB, id int) Setting {
	setting := Setting{}
	db.First(&setting, id)
	return setting
}

// FetchSettingByKey gets a setting from the database via the provided key
func FetchSettingByKey(db *gorm.DB, key string) Setting {
	setting := Setting{}
	db.Where("key = ?", key).First(&setting)
	return setting
}

// FetchSettings gets all settings from the database
func FetchSettings(db *gorm.DB) (settings []Setting) {
	db.Find(&settings)
	return settings
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

// DeleteSetting removes a setting with the matching id from the database
func DeleteSetting(db *gorm.DB, id int) Setting {
	setting := FetchSetting(db, id)
	db.Delete(&setting)
	return setting
}
