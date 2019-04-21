package cmd

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/markustenghamn/beubo/cmd/models"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var DB = setupDB()

func setupDB() *gorm.DB {
	db, err := gorm.Open("mysql", "root:h2BCZ39Q23@/beubo?charset=utf8&parseTime=True&loc=Local")
	checkErr(err)

	return db
}

func Init() {
	DB.AutoMigrate(
		&models.User{},
		&models.UserActivation{},
		&models.UserRole{},
		&models.Config{},
		&models.Post{},
		&models.Theme{},
		&models.Site{})
}

func Seed() {
	var err error
	email := "m@rkus.io"
	password := "Test1234!"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	checkErr(err)

	user := models.User{Email: email, Password: string(hashedPassword)}

	if DB.NewRecord(user) { // => returns `true` as primary key is blank
		DB.Create(&user)
	}
}
