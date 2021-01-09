package structs

import (
	"fmt"
	"github.com/uberswe/beubo/pkg/utility"
	"gorm.io/gorm"
)

// Role has one or more users and one or more features. A user belonging to a role which also has a feature will allow that user to use the feature
type Role struct {
	gorm.Model
	Name     string     `gorm:"size:255"`
	Users    []*User    `gorm:"many2many:user_roles;"`
	Features []*Feature `gorm:"many2many:role_features;"`
}

// Feature contains a key referencing features in the application
type Feature struct {
	gorm.Model
	Key   string  `gorm:"size:255"`
	Roles []*Role `gorm:"many2many:role_features;"`
}

func (r Role) IsDefault() bool {
	if r.Name == "Administrator" || r.Name == "Member" {
		return true
	}
	return false
}

func (r Role) HasFeature(db *gorm.DB, f Feature) bool {
	features := []Feature{}
	_ = db.Model(&r).Where("key = ?", f.Key).Association("Features").Find(&features)
	if len(features) >= 1 {
		return true
	}
	return false
}

// CreateRole is a method which creates a role using gorm
func CreateRole(db *gorm.DB, name string) bool {
	checkRole := FetchRoleByName(db, name)
	if checkRole.ID != 0 {
		return false
	}
	role := Role{Name: name}

	if err := db.Create(&role).Error; err != nil {
		fmt.Println("Could not create role")
		return false
	}
	return true
}

// FetchRoleByName retrieves a role from the database using the provided name
func FetchRoleByName(db *gorm.DB, name string) Role {
	role := Role{}
	db.Where("name = ?", name).First(&role)
	return role
}

// FetchRole retrieves a role from the database using the provided id
func FetchRole(db *gorm.DB, id int) Role {
	role := Role{}
	db.First(&role, id)
	return role
}

// DeleteRole deletes a user by id
func DeleteRole(db *gorm.DB, id int) (role Role) {
	db.First(&role, id)
	db.Delete(&role)
	return role
}

// UpdateRole updates the role struct with the provided details
func UpdateRole(db *gorm.DB, id int, name string, features []*Feature) bool {
	role := FetchRole(db, id)
	checkRole := FetchRoleByName(db, name)
	if checkRole.ID != 0 && checkRole.ID != role.ID {
		return false
	}
	err := db.Model(&role).Association("Features").Clear()
	utility.ErrorHandler(err, false)
	role.Name = name
	role.Features = features
	if err := db.Save(&role).Error; err != nil {
		fmt.Println("Could not create role")
		return false
	}
	return true
}
