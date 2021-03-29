package beubo

import (
	"fmt"
	"github.com/uberswe/beubo/pkg/middleware"
	"github.com/uberswe/beubo/pkg/plugin"
	"github.com/uberswe/beubo/pkg/structs"
	"github.com/uberswe/beubo/pkg/utility"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io/ioutil"
	"log"
	"os"
	"time"
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
	log.Println("Opening database...")
	dialector := getDialector(databaseUser, databasePassword, databaseHost, databasePort, databaseName, databaseDriver)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Silent, // Log level
			Colorful:      true,
		},
	)
	config := gorm.Config{
		Logger: newLogger,
	}
	if databaseDriver == "sqlite3" {
		config = gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   newLogger,
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

	log.Println("Running database migrations...")

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
		&structs.Feature{},
		&structs.MenuSection{},
		&structs.MenuItem{},
		&structs.MenuPermission{},
		&structs.MenuSetting{})
	utility.ErrorHandler(err, true)
}

func prepareSeed(email string, password string) {
	shouldSeed = true
	seedEmail = email
	seedPassword = password
}

func databaseSeed() {
	theme := addThemes()

	// user registration is disabled by default
	disableRegistration := structs.Setting{Key: "enable_user_registration", Value: "false"}
	DB.Where("key = ?", disableRegistration.Key).First(&disableRegistration)
	if disableRegistration.ID == 0 {
		DB.Create(&disableRegistration)
	}

	// users who register should have a member role
	newUserRole := structs.Setting{Key: "new_user_role", Value: "Member"}
	DB.Where("key = ?", newUserRole.Key).First(&newUserRole)
	if newUserRole.ID == 0 {
		DB.Create(&newUserRole)
	}

	features := []*structs.Feature{
		{Key: "manage_sites"},
		{Key: "manage_pages"},
		{Key: "manage_users"},
		{Key: "manage_user_roles"},
		{Key: "manage_plugins"},
		{Key: "manage_settings"},
		{Key: "manage_menus"},
	}

	for _, feature := range features {
		DB.Where("key = ?", feature.Key).First(&feature)
		if feature.ID == 0 {
			DB.Create(&feature)
		}
	}

	// Add default roles if not exist
	adminRole := structs.Role{}
	DB.Where("name = ?", "Administrator").First(&adminRole)
	if adminRole.ID == 0 {
		adminRole = structs.Role{Name: "Administrator", Features: features}
		DB.Create(&adminRole)
	}
	role := structs.Role{}
	DB.Where("name = ?", "Member").First(&role)
	if role.ID == 0 {
		role = structs.Role{Name: "Member"}
		DB.Create(&role)
	}

	menuItems := []structs.MenuItem{
		{
			Text:          "Home",
			URI:           "/",
			Authenticated: true,
		},
		{
			Text:          "Home",
			URI:           "/",
			Authenticated: false,
		},
		{
			Text:          "Admin",
			URI:           "/admin",
			Authenticated: true,
			Permissions: []structs.MenuPermission{
				{
					Permission: "manage_sites",
					Show:       true,
				},
			},
		},
		{
			Text:          "Login",
			URI:           "/login",
			Authenticated: false,
		},
		{
			Text:          "Logout",
			URI:           "/logout",
			Authenticated: true,
		},
		{
			Text: "Register",
			URI:  "/register",
			Settings: []structs.MenuSetting{
				{
					Setting:     "enable_user_registration",
					ShouldEqual: "true",
					Show:        true,
				},
			},
			Authenticated: false,
		},
	}

	var adminSidebarMenu []structs.MenuItem

	adminSidebarMenu = append(adminSidebarMenu, structs.MenuItem{
		Text:          "Sites",
		URI:           "/admin/",
		Authenticated: true,
		Permissions: []structs.MenuPermission{
			{
				Permission: "manage_sites",
				Show:       true,
			},
		},
	})

	adminSidebarMenu = append(adminSidebarMenu, structs.MenuItem{
		Text:          "Settings",
		URI:           "/admin/settings",
		Authenticated: true,
		Permissions: []structs.MenuPermission{
			{
				Permission: "manage_settings",
				Show:       true,
			},
		},
	})

	adminSidebarMenu = append(adminSidebarMenu, structs.MenuItem{
		Text:          "Users",
		URI:           "/admin/users",
		Authenticated: true,
		Permissions: []structs.MenuPermission{
			{
				Permission: "manage_users",
				Show:       true,
			},
		},
		Items: []structs.MenuItem{
			{
				Text:          "Roles",
				URI:           "/admin/users/roles",
				Authenticated: true,
				Permissions: []structs.MenuPermission{
					{
						Permission: "manage_user_roles",
						Show:       true,
					},
				},
			},
		},
	})

	adminSidebarMenu = append(adminSidebarMenu, structs.MenuItem{
		Text:          "Plugins",
		URI:           "/admin/plugins",
		Authenticated: true,
		Permissions: []structs.MenuPermission{
			{
				Permission: "manage_plugins",
				Show:       true,
			},
		},
	})

	adminSidebarMenu = append(adminSidebarMenu, structs.MenuItem{
		Text:          "Menus",
		URI:           "/admin/menus",
		Authenticated: true,
		Permissions: []structs.MenuPermission{
			{
				Permission: "manage_menus",
				Show:       true,
			},
		},
	})

	var adminSiteSidebarMenu []structs.MenuItem

	adminSiteSidebarMenu = append(adminSiteSidebarMenu, structs.MenuItem{
		Text:          "Pages",
		URI:           fmt.Sprintf("/admin/sites/a/{%s}/", middleware.AdminSiteContextKey),
		Authenticated: true,
		Permissions: []structs.MenuPermission{
			{
				Permission: "manage_pages",
				Show:       true,
			},
		},
	})
	adminSiteSidebarMenu = append(adminSiteSidebarMenu, structs.MenuItem{
		Text:          "Plugins",
		URI:           fmt.Sprintf("/admin/sites/a/{%s}/plugins", middleware.AdminSiteContextKey),
		Authenticated: true,
	})
	adminSiteSidebarMenu = append(adminSiteSidebarMenu, structs.MenuItem{
		Text:          "Menus",
		URI:           fmt.Sprintf("/admin/sites/a/{%s}/menus", middleware.AdminSiteContextKey),
		Authenticated: true,
	})

	tmp := make([]structs.MenuItem, len(menuItems))
	copy(tmp, menuItems)

	menus := []structs.MenuSection{
		{
			Items:   menuItems,
			Section: "header",
			SiteID:  0,
		},
		{
			Items:    adminSidebarMenu,
			Section:  "admin_sidebar",
			SiteID:   0,
			Template: "menu.sidebar",
		},
		{
			Items:    adminSiteSidebarMenu,
			Section:  "admin_site_sidebar",
			SiteID:   1,
			Template: "menu.sidebar",
		},
		{
			Items:   tmp,
			Section: "header",
			SiteID:  1,
		},
	}
	for _, menu := range menus {
		DB.Where("section = ?", menu.Section).Where("site_id", menu.SiteID).First(&menu)
		if menu.ID == 0 {
			DB.Create(&menu)
		}
	}

	// Add the specified default test user if the environment is also not set to production
	if environment != "production" && testuser != "" && testpass != "" {
		user := structs.User{Email: testuser}
		DB.Where("email = ?", user.Email).First(&user)
		if user.ID == 0 {
			// ASVS 4.0 point 2.4.4 states cost should be at least 13 https://github.com/OWASP/ASVS/
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testpass), 14)
			utility.ErrorHandler(err, true)
			user.Password = string(hashedPassword)
			user.Roles = []*structs.Role{
				&adminRole,
			}
			DB.Create(&user)
		}
	}

	log.Println("should seed", shouldSeed)

	// If seeding is enabled we perform the seed with default info
	if shouldSeed {
		seedData(theme, adminRole)
	}

	shouldSeed = false
	seedEmail = ""
	seedPassword = ""
}

func seedData(theme structs.Theme, adminRole structs.Role) {
	log.Println("Seeding database...")
	user := structs.User{Email: seedEmail}
	DB.Where("email = ?", user.Email).First(&user)
	if user.ID == 0 {
		// ASVS 4.0 point 2.4.4 states cost should be at least 13 https://github.com/OWASP/ASVS/
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(seedPassword), 14)
		user.Password = string(hashedPassword)

		utility.ErrorHandler(err, true)
		user.Roles = []*structs.Role{
			&adminRole,
		}
		DB.Create(&user)
	}

	// Create a site

	site := structs.Site{
		Title:  "Beubo",
		Domain: "beubo.localhost",
		Theme:  theme,
		Type:   1,
	}

	DB.Where("domain = ?", site.Domain).First(&site)
	if site.ID == 0 {
		DB.Create(&site)
		// Seed permissions
		var users []structs.User
		DB.Preload("Sites").Find(&users)
		for _, user := range users {
			exists := false
			for _, s := range user.Sites {
				if s.ID == site.ID {
					exists = true
					break
				}
			}
			if !exists {
				user.Sites = append(user.Sites, &site)
				DB.Save(&user)
			}
		}
	}

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
	DB.Where("slug = ?", page.Slug).First(&page)
	if page.ID == 0 {
		DB.Create(&page)
	}
}

func addThemes() (theme structs.Theme) {
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
	return theme
}
