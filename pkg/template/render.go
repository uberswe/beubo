package template

import (
	"fmt"
	"github.com/uberswe/beubo/pkg/middleware"
	"github.com/uberswe/beubo/pkg/plugin"
	"github.com/uberswe/beubo/pkg/structs"
	"github.com/uberswe/beubo/pkg/structs/page"
	"github.com/uberswe/beubo/pkg/structs/page/menu"
	"github.com/uberswe/beubo/pkg/utility"
	"gorm.io/gorm"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	currentTheme = "default"
	rootDir      = "./"
)

// BeuboTemplateRenderer holds all the configuration variables for rendering templates in Beubo
type BeuboTemplateRenderer struct {
	T               *template.Template
	PluginHandler   *plugin.Handler
	ReloadTemplates bool
	CurrentTheme    string
	ThemeDir        string
	DB              *gorm.DB
}

// Init prepares the BeuboTemplateRenderer to render pages with html templates
func (btr *BeuboTemplateRenderer) Init() {
	log.Println("Parsing and loading templates...")
	var err error
	btr.T, err = findAndParseTemplates(rootDir)
	utility.ErrorHandler(err, false)
}

// RenderHTMLPage handles rendering of the html template and should be the last function called before returning the response
func (btr *BeuboTemplateRenderer) RenderHTMLPage(w http.ResponseWriter, r *http.Request, pageData structs.PageData) {

	if os.Getenv("ASSETS_DIR") != "" {
		rootDir = os.Getenv("ASSETS_DIR")
	}

	// Session flash messages to prompt failed logins etc..
	errorMessage, err := utility.GetFlash(w, r, "error")
	utility.ErrorHandler(err, false)
	warningMessage, err := utility.GetFlash(w, r, "warning")
	utility.ErrorHandler(err, false)
	stringMessage, err := utility.GetFlash(w, r, "message")
	utility.ErrorHandler(err, false)

	siteName := "Beubo"

	// Get the site from context
	site := r.Context().Value(middleware.SiteContextKey)

	if site != nil {
		btr.CurrentTheme = site.(structs.Site).Theme.Slug
		siteName = site.(structs.Site).Title
	} else if os.Getenv("THEME") != "" {
		btr.CurrentTheme = os.Getenv("THEME")
	} else {
		// Default theme
		if btr.CurrentTheme == "" {
			btr.CurrentTheme = "default"
		}
	}

	// TODO in the future we should make some way for the theme to define the stylesheets
	scripts := []string{
		"/default/js/main.js",
	}

	stylesheets := []string{
		"/default/css/main.min.css",
		"/default/css/milligram-1-4-1.min.css",
		"/default/css/normalize-8-0-1.min.css",
	}

	if strings.HasPrefix(r.URL.Path, "/admin") {
		scripts = append(scripts, "/default/js/admin.js")
		stylesheets = append(stylesheets, "/default/css/admin.css")
	}

	data := structs.PageData{
		Stylesheets: stylesheets,
		// TODO make the favicon dynamic
		Favicon:     "/default/images/favicon.ico",
		Scripts:     scripts,
		WebsiteName: siteName,
		URL:         "http://localhost:3000",
		Error:       string(errorMessage),
		Warning:     string(warningMessage),
		Message:     string(stringMessage),
		Year:        strconv.Itoa(time.Now().Year()),
		Menus:       btr.buildMenus(r),
	}

	data = mergePageData(data, pageData)
	// PluginHandler is not available during installation
	if btr.PluginHandler != nil {
		data = btr.PluginHandler.PageData(r, data)
	}

	if btr.ReloadTemplates {
		log.Println("Parsing and loading templates...")
		var err error
		btr.T, err = findAndParseTemplates(rootDir)
		utility.ErrorHandler(err, false)
	}

	var foundTemplate *template.Template

	path := fmt.Sprintf("%s.%s", btr.CurrentTheme, pageData.Template)
	if foundTemplate = btr.T.Lookup(path); foundTemplate == nil {
		// Fallback to default
		path := fmt.Sprintf("%s.%s", "default", pageData.Template)
		if foundTemplate = btr.T.Lookup(path); foundTemplate == nil {
			log.Printf("Theme file not found %s\n", path)
			return
		}
	}

	err = foundTemplate.Execute(w, data)
	utility.ErrorHandler(err, false)
}

func (btr *BeuboTemplateRenderer) buildMenus(r *http.Request) []page.Menu {
	userFromContext := r.Context().Value(middleware.UserContextKey)
	adminSiteFromContext := r.Context().Value(middleware.AdminSiteContextKey)
	var user structs.User
	if userFromContext != nil && userFromContext.(structs.User).ID > 0 {
		user = userFromContext.(structs.User)
	}
	var adminSite structs.Site
	if adminSiteFromContext != nil && adminSiteFromContext.(structs.Site).ID > 0 {
		adminSite = adminSiteFromContext.(structs.Site)
	}

	// TODO query database and fetch menu config
	menuItems := []page.MenuItem{
		{Text: "Home", URI: "/"},
		{Text: "Login", URI: "/login"},
	}

	// DB is not available during installation
	if btr.DB != nil {
		setting := structs.FetchSettingByKey(btr.DB, "enable_user_registration")
		if !(setting.ID == 0 || setting.Value == "false") {
			menuItems = append(menuItems, page.MenuItem{Text: "Register", URI: "/register"})
		}
	}

	menus := []page.Menu{menu.DefaultMenu{
		Items:      menuItems,
		Identifier: "header",
		T:          btr.T,
	}}

	if user.ID > 0 {
		adminHeaderMenu := []page.MenuItem{
			{Text: "Home", URI: "/"},
			{Text: "Logout", URI: "/logout"},
		}

		if user.CanAccess(btr.DB, "manage_sites") ||
			user.CanAccess(btr.DB, "manage_pages") ||
			user.CanAccess(btr.DB, "manage_users") ||
			user.CanAccess(btr.DB, "manage_user_roles") ||
			user.CanAccess(btr.DB, "manage_plugins") ||
			user.CanAccess(btr.DB, "manage_settings") {
			adminHeaderMenu = []page.MenuItem{
				{Text: "Home", URI: "/"},
				{Text: "Admin", URI: "/admin"},
				{Text: "Logout", URI: "/logout"},
			}
		}

		adminSidebarMenu := []page.MenuItem{}

		if user.CanAccess(btr.DB, "manage_sites") {
			adminSidebarMenu = append(adminSidebarMenu, page.MenuItem{Text: "Sites", URI: "/admin/"})
		}

		if user.CanAccess(btr.DB, "manage_settings") {
			adminSidebarMenu = append(adminSidebarMenu, page.MenuItem{Text: "Settings", URI: "/admin/settings"})
		}

		if user.CanAccess(btr.DB, "manage_users") {
			adminSidebarMenu = append(adminSidebarMenu, page.MenuItem{
				Text: "Users",
				URI:  "/admin/users",
				Items: []page.MenuItem{
					{
						Text: "Roles",
						URI:  "/admin/users/roles",
					},
				},
				// Submenus need template to be defined
				T: btr.T,
			})
		}

		if user.CanAccess(btr.DB, "manage_plugins") {
			adminSidebarMenu = append(adminSidebarMenu, page.MenuItem{Text: "Plugins", URI: "/admin/plugins"})
		}

		adminSiteSidebarMenu := []page.MenuItem{}

		if user.CanAccess(btr.DB, "manage_pages") && adminSite.ID > 0 {
			// TODO check if user has permission to access adminSite
			adminSiteSidebarMenu = append(adminSiteSidebarMenu, page.MenuItem{Text: "Pages", URI: fmt.Sprintf("/admin/sites/a/%d/", adminSite.ID)})
			adminSiteSidebarMenu = append(adminSiteSidebarMenu, page.MenuItem{Text: "Plugins", URI: fmt.Sprintf("/admin/sites/a/%d/plugins", adminSite.ID)})
		}

		menus = []page.Menu{menu.DefaultMenu{
			Items:      adminHeaderMenu,
			Identifier: "header",
			T:          btr.T,
		}, menu.DefaultMenu{
			Items:      adminSidebarMenu,
			Identifier: "admin_sidebar",
			Template:   "menu.sidebar",
			T:          btr.T,
		}, menu.DefaultMenu{
			Items:      adminSiteSidebarMenu,
			Identifier: "admin_site_sidebar",
			Template:   "menu.sidebar",
			T:          btr.T,
		}}
	}
	return menus
}

// findAndParseTemplates finds all the templates in the rootDir and makes a template map
// This method was found here https://stackoverflow.com/a/50581032/1260548
func findAndParseTemplates(rootDir string) (*template.Template, error) {
	cleanRoot, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, err
	}
	pfx := len(cleanRoot) + 1
	root := template.New("")

	err = filepath.Walk(cleanRoot, func(path string, info os.FileInfo, e1 error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			if e1 != nil {
				return e1
			}

			b, e2 := ioutil.ReadFile(path)
			if e2 != nil {
				return e2
			}

			name := path[pfx:]
			t := root.New(name)
			t, e2 = t.Parse(string(b))
			if e2 != nil {
				return e2
			}
		}

		return nil
	})

	return root, err
}

func mergePageData(a structs.PageData, b structs.PageData) structs.PageData {
	// TODO this could be simplified by making a function that compares an interface and picks a value but I decided that this is more readable for now
	if b.Template != "" {
		a.Template = b.Template
	}

	if b.Title != "" {
		a.Title = b.Title
	}

	if b.WebsiteName != "" {
		a.WebsiteName = b.WebsiteName
	}

	if b.Error != "" {
		a.Error = b.Error
	}

	if b.Warning != "" {
		a.Warning = b.Warning
	}

	if b.Message != "" {
		a.Message = b.Message
	}

	if len(b.Scripts) > 0 {
		a.Scripts = b.Scripts
	}

	if len(b.Stylesheets) > 0 {
		a.Stylesheets = b.Stylesheets
	}

	if b.Favicon != "" {
		a.Favicon = b.Favicon
	}

	if len(b.Menus) > 0 {
		a.Menus = b.Menus
	}

	a.Components = b.Components
	a.Themes = b.Themes
	a.Templates = b.Templates
	a.Extra = b.Extra

	return a
}
