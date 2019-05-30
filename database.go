package beubo

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/markustenghamn/beubo/pkg/models"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var (
	seedEmail    = ""
	seedPassword = ""
	shouldSeed   = false
	DB           *gorm.DB
)

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

func prepareSeed(email string, password string) {
	shouldSeed = true
	seedEmail = email
	seedPassword = password
}

func databaseSeed() {
	if shouldSeed {
		var err error

		// ASVS 4.0 point 2.4.4 states cost should be at least 13 https://github.com/OWASP/ASVS/
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(seedPassword), 14)

		checkErr(err)

		user := models.User{Email: seedEmail, Password: string(hashedPassword)}

		if DB.NewRecord(user) { // => returns `true` as primary key is blank
			DB.Create(&user)
		}
		shouldSeed = false
		seedEmail = ""
		seedPassword = ""
	}
}
