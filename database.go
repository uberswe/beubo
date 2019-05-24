package beubo

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/markustenghamn/beubo/pkg/models"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var DB *gorm.DB

func setupDB() *gorm.DB {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", databaseUser, databasePassword, databaseHost, databasePort, databaseName))
	checkErr(err)

	return db
}

func databaseInit() {
	DB = setupDB()

	type Result struct {
		DropQuery string
	}

	var result []Result

	DB.Raw("SELECT concat('DROP TABLE IF EXISTS `', table_name, '`;') as drop_query FROM information_schema.tables WHERE table_schema = 'beubo';").Scan(&result)

	for _, r := range result {
		DB.Exec(r.DropQuery)
	}

	DB.AutoMigrate(
		&models.User{},
		&models.UserActivation{},
		&models.UserRole{},
		&models.Config{},
		&models.Post{},
		&models.Theme{},
		&models.Site{})
}

func databaseSeed() {
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
