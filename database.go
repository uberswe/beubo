package beubo

import (
	"fmt"
	"github.com/uberswe/beubo/pkg/plugin"
	"github.com/uberswe/beubo/pkg/utility"
	"gorm.io/gorm"
	"io/ioutil"
	"log"

	"github.com/uberswe/beubo/pkg/structs"
	"golang.org/x/crypto/bcrypt"
	// Gorm recommends a blank import to support underlying mysql
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
)

var (
	seedEmail    = "seed@beubo.com"
	seedPassword = "Beubo1234!"
	// TODO change this to a config
	shouldSeed            = false
	shouldRefreshDatabase = false
	// DB is used to perform database queries globally. In the future this should probably
	// be changed so that database.go declares methods that can be used to perform types of
	// queries
	DB *gorm.DB
)

func setupDB() *gorm.DB {
	log.Println("Opening database")
	dialector := getDialector(databaseUser, databasePassword, databaseHost, databasePort, databaseName, databaseDriver)
	config := gorm.Config{}
	if databaseDriver == "sqlite3" {
		config = gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		}
	}
	db, err := gorm.Open(dialector, &config)
	utility.ErrorHandler(err, true)

	return db
}

func getDialector(user string, pass string, host string, port string, name string, driver string) gorm.Dialector {
	connectString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, pass, host, port, name)
	dialector := mysql.Open(connectString)
	if driver == "sqlite3" {
		connectString = databaseName
		dialector = sqlite.Open(connectString)
	}
	return dialector
}

func databaseInit() {
	DB = setupDB()

	if shouldRefreshDatabase {

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
	}

	log.Println("Running database migrations")

	err := DB.AutoMigrate(
		&structs.User{},
		&structs.UserActivation{},
		&structs.Config{},
		&structs.Page{},
		&structs.Theme{},
		&structs.Session{},
		&structs.Site{},
		&structs.Tag{},
		&structs.Comment{},
		&structs.Setting{},
		&plugin.PluginSite{},
		&structs.Role{},
		&structs.Feature{})
	utility.ErrorHandler(err, true)
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
			theme = structs.Theme{}
			DB.Where("slug = ?", file.Name()).First(&theme)
			if theme.ID == 0 {
				theme = structs.Theme{Slug: file.Name(), Title: file.Name()}
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
		DB.Where("email = ?", user.Email).First(&user)
		if user.ID == 0 {
			DB.Create(&user)
		}

	}

	// user registration is disabled by default
	disableRegistration := structs.Setting{Key: "enable_user_registration", Value: "false"}
	DB.Where("key = ?", disableRegistration.Key).First(&disableRegistration)
	if disableRegistration.ID == 0 {
		DB.Create(&disableRegistration)
	}

	// If seeding is enabled we perform the seed with default info
	if shouldSeed {
		log.Println("Seeding database")
		var err error

		// ASVS 4.0 point 2.4.4 states cost should be at least 13 https://github.com/OWASP/ASVS/
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(seedPassword), 14)

		utility.ErrorHandler(err, true)

		user := structs.User{Email: seedEmail, Password: string(hashedPassword)}

		DB.Create(&user)

		// Create a site

		site := structs.Site{
			Title:  "Beubo",
			Domain: "beubo.localhost",
			Theme:  theme,
			Type:   1,
		}

		DB.Create(&site)

		// Create a page
		content := `<p>Welcome to Beubo! Beubo is a free, simple, and minimal CMS with unlimited extensibility using plugins. This is the default page and can be changed in the admin area for this site.</p>`
		content += `<p>Beubo is open source and the project can be found on <a href="https://github.com/uberswe/beubo">Github</a>. If you find any problems or have an idea on how Beubo can be improved, please feel free to <a href="https://github.com/uberswe/beubo/issues">open an issue here</a>.</p>`
		content += `<p>Feel free to <a href="https://github.com/uberswe/beubo/pulls">open a pull request</a> if you would like to contribute your own changes.</p>`
		content += `<p>For more information on how to use, customize and extend Beubo please see the <a href="https://github.com/uberswe/beubo/wiki">wiki</a></p>`

		page := structs.Page{
			Model:    gorm.Model{},
			Title:    "Default page",
			Content:  content,
			Slug:     "/",
			Template: "page",
			SiteID:   int(site.ID),
		}

		DB.Create(&page)

		shouldSeed = false
		seedEmail = ""
		seedPassword = ""
	}
}
