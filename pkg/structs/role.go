package structs

import "gorm.io/gorm"

// Role has one or more users and one or more features. A user belonging to a role which also has a feature will allow that user to use the feature
type Role struct {
	gorm.Model
	Name     string     `gorm:"size:255"`
	Users    []*User    `gorm:"many2many:user_roles;"`
	Features []*Feature `gorm:"many2many:role_features;"`
}

type Feature struct {
	gorm.Model
	Key   string  `gorm:"size:255"`
	Roles []*Role `gorm:"many2many:role_features;"`
}
