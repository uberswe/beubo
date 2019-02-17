package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Email       string `gorm:"size:255;unique_index"`
	Password    string `gorm:"size:255"`
	Activations []UserActivation
	Roles       []UserRole
}

type UserActivation struct {
	gorm.Model
	UserID uint
	Type   string // Email, SMS, PushNotification
	Active bool
	Code   string
}

type UserRole struct {
	gorm.Model
	UserID uint
	Name   string
}

func CreateUser(db *gorm.DB, email string, password string) bool {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return false
	}

	user := User{Email: email, Password: string(hashedPassword)}

	if db.NewRecord(user) { // => returns `true` as primary key is blank
		db.Create(&user)
	}
	return true
}

func AuthUser(db *gorm.DB, email string, password string) bool {
	var user User

	if db.Model(&user).Where("email = ?", email).RecordNotFound() {
		return false
	}

	db.Where("email = ?", email).First(&user)

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false
	} else if err != nil {
		return false
	}

	return true
}
