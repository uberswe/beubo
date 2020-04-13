package beubo

import (
	"fmt"
	"github.com/jinzhu/gorm"
	// Gorm recommends a blank import to support underlying mysql
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/markustenghamn/beubo/pkg/structs"
	"golang.org/x/crypto/bcrypt"
)

var (
	seedEmail    = ""
	seedPassword = ""
	shouldSeed   = false
	// DB is used to perform database queries globally. In the future this should probably
	// be changed so that database.go declares methods that can be used to perform types of
	// queries
	DB *gorm.DB
)

func setupDB() *gorm.DB {
	connectString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", databaseUser, databasePassword, databaseHost, databasePort, databaseName)
	if databaseDriver == "sqlite3" {
		connectString = "beubo.db"
	}
	db, err := gorm.Open(databaseDriver, connectString)
	checkErr(err)

	return db
}

func databaseInit() {
	DB = setupDB()

	type Result struct {
		DropQuery string
	}

	if databaseDriver != "sqlite3" {

		var result []Result
		// TODO if sqlite3 we can not use information_schema here
		DB.Raw("SELECT concat('DROP TABLE IF EXISTS `', table_name, '`;') as drop_query FROM information_schema.tables WHERE table_schema = 'beubo';").Scan(&result)

		for _, r := range result {
			DB.Exec(r.DropQuery)
		}

	}

	DB.AutoMigrate(
		&structs.User{},
		&structs.UserActivation{},
		&structs.UserRole{},
		&structs.Config{},
		&structs.Page{},
		&structs.Theme{},
		&structs.Site{})
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

		user := structs.User{Email: seedEmail, Password: string(hashedPassword)}

		if DB.NewRecord(user) { // => returns `true` as primary key is blank
			DB.Create(&user)
		}
		shouldSeed = false
		seedEmail = ""
		seedPassword = ""
	}
}
