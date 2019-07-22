package models

import (
	"fmt"
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

	user := User{Email: email}

	if db.NewRecord(user) { // => returns `true` as primary key is blank
		if err := db.Create(&user).Error; err != nil {
			fmt.Println("Could not create user")
			return false
		}
		// User password is hashed after the response is returned to improve performance
		go hashUserPassword(db, user, password)
		return true
	}
	return false
}

func hashUserPassword(db *gorm.DB, user User, password string) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		fmt.Println("Password hashing failed for user")
		return
	}

	user.Password = string(hashedPassword)

	if err := db.Save(&user).Error; err != nil {
		fmt.Println("Could not update hashed password for user")
	}
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
