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
	user := r.Context().Value(middleware.UserContextKey)

	if site != nil {
		btr.CurrentTheme = site.(structs.Site).Theme.Slug
		siteName = site.(structs.Site).Title
	} else if os.Getenv("THEME") != "" {
		btr.CurrentTheme = os.Getenv("THEME")
	} else {
		// Default theme
		btr.CurrentTheme = "default"
	}

	menuItems := []page.MenuItem{
		{Text: "Home", URI: "/"},
		{Text: "Login", URI: "/login"},
	}

	setting := structs.FetchSettingByKey(btr.DB, "enable_user_registration")
	if !(setting.ID == 0 || setting.Value == "false") {
		menuItems = append(menuItems, page.MenuItem{Text: "Register", URI: "/register"})
	}

	menus := []page.Menu{menu.DefaultMenu{
		Items:      menuItems,
		Identifier: "header",
		T:          btr.T,
	}}

	if user != nil && user.(structs.User).ID > 0 {
		adminHeaderMenu := []page.MenuItem{
			{Text: "Home", URI: "/"},
			{Text: "Logout", URI: "/logout"},
		}

		if user.(structs.User).CanAccess(btr.DB, "manage_sites") ||
			user.(structs.User).CanAccess(btr.DB, "manage_pages") ||
			user.(structs.User).CanAccess(btr.DB, "manage_users") ||
			user.(structs.User).CanAccess(btr.DB, "manage_user_roles") ||
			user.(structs.User).CanAccess(btr.DB, "manage_plugins") ||
			user.(structs.User).CanAccess(btr.DB, "manage_settings") {
			adminHeaderMenu = []page.MenuItem{
				{Text: "Home", URI: "/"},
				{Text: "Admin", URI: "/admin"},
				{Text: "Logout", URI: "/logout"},
			}
		}

		adminSidebarMenu := []page.MenuItem{}

		if user.(structs.User).CanAccess(btr.DB, "manage_sites") {
			adminSidebarMenu = append(adminSidebarMenu, page.MenuItem{Text: "Sites", URI: "/admin/"})
		}

		if user.(structs.User).CanAccess(btr.DB, "manage_settings") {
			adminSidebarMenu = append(adminSidebarMenu, page.MenuItem{Text: "Settings", URI: "/admin/settings"})
		}

		if user.(structs.User).CanAccess(btr.DB, "manage_users") {
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

		if user.(structs.User).CanAccess(btr.DB, "manage_plugins") {
			adminSidebarMenu = append(adminSidebarMenu, page.MenuItem{Text: "Plugins", URI: "/admin/plugins"})
		}

		menus = []page.Menu{menu.DefaultMenu{
			Items:      adminHeaderMenu,
			Identifier: "header",
			T:          btr.T,
		}, menu.DefaultMenu{
			Items:      adminSidebarMenu,
			Identifier: "sidebar",
			Template:   "menu.sidebar",
			T:          btr.T,
		}}
	}

	// TODO in the future we should make some way for the theme to define the stylesheets
	scripts := []string{
		"/default/js/main.js",
	}

	stylesheets := []string{
		"/default/css/main.css",
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
		Menus:       menus,
	}

	data = mergePageData(data, pageData)
	data = btr.PluginHandler.PageData(r, data)

	if btr.ReloadTemplates {
		log.Println("Parsing and loading templates...")
		var err error
		btr.T, err = findAndParseTemplates(rootDir)
		utility.ErrorHandler(err, false)
	}

	var foundTemplate *template.Template

	path := fmt.Sprintf("%s.%s", currentTheme, pageData.Template)
	if foundTemplate = btr.T.Lookup(path); foundTemplate == nil {
		log.Printf("Theme file not found %s\n", path)
		return
	}

	err = foundTemplate.Execute(w, data)
	utility.ErrorHandler(err, false)
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

	a.Components = b.Components
	a.Themes = b.Themes
	a.Templates = b.Templates
	a.Extra = b.Extra

	return a
}
