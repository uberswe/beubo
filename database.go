package beubo

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/markustenghamn/beubo/pkg/utility"
	"io/ioutil"
	"log"

	// Gorm recommends a blank import to support underlying mysql
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/markustenghamn/beubo/pkg/structs"
	"golang.org/x/crypto/bcrypt"
)

var (
	seedEmail    = ""
	seedPassword = ""
	// TODO change this to a config
	shouldSeed = true
	// DB is used to perform database queries globally. In the future this should probably
	// be changed so that database.go declares methods that can be used to perform types of
	// queries
	DB *gorm.DB
)

func setupDB() *gorm.DB {
	log.Println("Opening database")
	connectString := ""
	if databaseDriver == "sqlite3" {
		connectString = databaseName
	} else {
		connectString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", databaseUser, databasePassword, databaseHost, databasePort, databaseName)
	}
	db, err := gorm.Open(databaseDriver, connectString)
	utility.ErrorHandler(err, true)

	return db
}

func databaseInit() {
	DB = setupDB()

	type Result struct {
		DropQuery string
	}

	var result []Result

	log.Println("Dropping all database tables")

	if databaseDriver == "sqlite3" {
		DB.Raw("SELECT 'DROP TABLE IF EXISTS `' || name || '`;'  as drop_query FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%';").Scan(&result)
	} else {
		DB.Raw("SELECT concat('DROP TABLE IF EXISTS `', table_name, '`;') as drop_query FROM information_schema.tables WHERE table_schema = 'beubo';").Scan(&result)
	}

	for _, r := range result {
		DB.Exec(r.DropQuery)
	}

	log.Println("Running database migrations")

	DB.AutoMigrate(
		&structs.User{},
		&structs.UserActivation{},
		&structs.UserRole{},
		&structs.Config{},
		&structs.Page{},
		&structs.Theme{},
		&structs.Session{},
		&structs.Site{},
		&structs.Tag{},
		&structs.Comment{})
}

func prepareSeed(email string, password string) {
	shouldSeed = true
	seedEmail = email
	seedPassword = password
}

func databaseSeed() {

	theme := structs.Theme{}
	// Add initial themes
	files, err := ioutil.ReadDir(rootDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		// Ignore the install directory, only used for installation
		if file.IsDir() && file.Name() != "install" {
			theme = structs.Theme{Slug: file.Name(), Title: file.Name()}
			if DB.NewRecord(theme) { // => returns `true` as primary key is blank
				DB.Create(&theme)
			}
		}
	}

	// Add the specified default test user if the environment is also not set to production
	if environment != "production" && testuser != "" && testpass != "" {
		var err error

		// ASVS 4.0 point 2.4.4 states cost should be at least 13 https://github.com/OWASP/ASVS/
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testpass), 14)

		utility.ErrorHandler(err, true)

		user := structs.User{Email: testuser, Password: string(hashedPassword)}

		if DB.NewRecord(user) { // => returns `true` as primary key is blank
			DB.Create(&user)
		}
	}

	// If seeding is enabled we perform the seed with default info
	if shouldSeed {
		log.Println("Seeding database")
		var err error

		// ASVS 4.0 point 2.4.4 states cost should be at least 13 https://github.com/OWASP/ASVS/
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(seedPassword), 14)

		utility.ErrorHandler(err, true)

		user := structs.User{Email: seedEmail, Password: string(hashedPassword)}

		if DB.NewRecord(user) { // => returns `true` as primary key is blank
			DB.Create(&user)
		}

		// Create a site

		site := structs.Site{
			Title:  "Default",
			Domain: "localhost:3000",
			Theme:  theme,
			Type:   1,
		}

		if DB.NewRecord(site) {
			DB.Create(&site)
		}

		// Create a page

		content := `<p>This is a default page</p>`

		page := structs.Page{
			Model:    gorm.Model{},
			Title:    "Default page",
			Content:  content,
			Slug:     "/",
			Template: "page",
			SiteID:   int(site.ID),
		}

		if DB.NewRecord(page) {
			DB.Create(&page)
		}

		shouldSeed = false
		seedEmail = ""
		seedPassword = ""
	}
}
