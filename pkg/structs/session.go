package structs

import (
	"fmt"
	"github.com/uberswe/beubo/pkg/utility"
	"gorm.io/gorm"
)

// Session represents an authenticated user session, there can be multiple sessions for one user
type Session struct {
	gorm.Model
	Token  string `gorm:"size:255;unique_index"`
	UserID int
	User   User
}

// CreateSession is a method which creates a session using gorm
func CreateSession(db *gorm.DB, userID int) Session {
	token, err := utility.GenerateToken(255)
	utility.ErrorHandler(err, false)

	session := Session{
		Token:  token,
		UserID: userID,
	}

	if err := db.Create(&session).Error; err != nil {
		fmt.Println("Could not create session")
	}

	return session
}

// FetchUserFromSession takes a provided token string and fetches the user for the session matching the provided token
func FetchUserFromSession(db *gorm.DB, token string) User {
	user := User{}
	session := Session{}

	db.Where("token = ?", token).First(&session)

	if session.ID != 0 {
		user = FetchUser(db, session.UserID)
	}

	return user
}
