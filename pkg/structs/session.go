package structs

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/markustenghamn/beubo/pkg/utility"
)

// User is a user who can authenticate with Beubo
type Session struct {
	gorm.Model
	Token  string `gorm:"size:255;unique_index"`
	UserID int
	User   User
}

// CreateUser is a method which creates a user using gorm
func CreateSession(db *gorm.DB, userID int) Session {
	token, err := utility.GenerateToken(255)
	utility.ErrorHandler(err, false)

	session := Session{
		Token:  token,
		UserID: userID,
	}

	if db.NewRecord(session) { // => returns `true` as primary key is blank
		if err := db.Create(&session).Error; err != nil {
			fmt.Println("Could not create session")
		}
	}
	return session
}

func FetchUserFromSession(db *gorm.DB, token string) User {
	user := User{}
	session := Session{}

	db.Where("token = ?", token).First(&session)

	if session.ID != 0 {
		user = FetchUser(db, session.UserID)
	}

	return user
}
