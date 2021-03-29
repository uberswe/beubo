package template

import (
	"fmt"
	"github.com/uberswe/beubo/pkg/middleware"
	"github.com/uberswe/beubo/pkg/plugin"
	"github.com/uberswe/beubo/pkg/structs"
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

func (btr *BeuboTemplateRenderer) buildMenus(r *http.Request) []structs.MenuSection {
	siteID := 0
	userFromContext := r.Context().Value(middleware.UserContextKey)
	adminSiteFromContext := r.Context().Value(middleware.AdminSiteContextKey)
	siteFromContext := r.Context().Value(middleware.SiteContextKey)
	var user structs.User
	if userFromContext != nil && userFromContext.(structs.User).ID > 0 {
		user = userFromContext.(structs.User)
	}
	var Site structs.Site
	if siteFromContext != nil && siteFromContext.(structs.Site).ID > 0 {
		Site = siteFromContext.(structs.Site)
		// If we are in the admin area we are not viewing a specific site and want global menus
		if !strings.HasPrefix(r.URL.Path, "/admin") {
			siteID = int(Site.ID)
		}
	}
	var adminSite structs.Site
	if adminSiteFromContext != nil && adminSiteFromContext.(structs.Site).ID > 0 {
		adminSite = adminSiteFromContext.(structs.Site)
		siteID = int(adminSite.ID)
	}

	menus := structs.FetchMenusBySiteID(btr.DB, siteID)

	for i, section := range menus {
		menus[i].T = btr.T
		menus[i].Items = setRenderForMenuItems(btr, structs.FetchMenuItemsBySectionId(btr.DB, int(section.ID)), &user, siteID)
	}

	return menus
}

func setRenderForMenuItems(btr *BeuboTemplateRenderer, items []structs.MenuItem, user *structs.User, siteID int) []structs.MenuItem {
	var result []structs.MenuItem
	for _, item := range items {
		// Parse and replace url parameters
		if strings.Contains(item.URI, fmt.Sprintf("{%s}", middleware.SiteContextKey)) {
			item.URI = strings.Replace(item.URI, fmt.Sprintf("{%s}", middleware.SiteContextKey), fmt.Sprintf("%d", siteID), -1)
		}
		if strings.Contains(item.URI, fmt.Sprintf("{%s}", middleware.AdminSiteContextKey)) {
			item.URI = strings.Replace(item.URI, fmt.Sprintf("{%s}", middleware.AdminSiteContextKey), fmt.Sprintf("%d", siteID), -1)
		}
		if strings.Contains(item.URI, fmt.Sprintf("{%s}", middleware.UserContextKey)) {
			item.URI = strings.Replace(item.URI, fmt.Sprintf("{%s}", middleware.UserContextKey), fmt.Sprintf("%d", user.ID), -1)
		}
		err := btr.DB.Model(&item).Association("Permissions").Find(&item.Permissions)
		utility.ErrorHandler(err, false)
		for _, p := range item.Permissions {
			if !user.CanAccess(btr.DB, p.Permission) {
				// If the user doesn't have permission then we continue to the next item
				continue
			}
		}
		err = btr.DB.Model(&item).Association("Settings").Find(&item.Settings)
		utility.ErrorHandler(err, false)
		denied := false
		for _, s := range item.Settings {
			setting := structs.FetchSettingByKey(btr.DB, s.Setting)
			if setting.ID > 0 {
				if setting.Value != s.ShouldEqual {
					if s.Show {
						// If the settings don't match we go on to the next item
						denied = true
						continue
					}
				} else {
					if !s.Show {
						denied = true
						continue
					}
				}
			} else {
				if s.Show {
					denied = true
					continue
				}
			}
		}
		if denied {
			continue
		}
		// Process submenu items
		item.T = btr.T
		item.Items = setRenderForMenuItems(btr, structs.FetchMenuItemsByParentId(btr.DB, int(item.ID)), user, siteID)
		// Check for authentication
		if item.Authenticated && user != nil && user.ID > 0 {
			result = append(result, item)
		} else if !item.Authenticated && (user == nil || user.ID == 0) {
			result = append(result, item)
		}
	}
	return result
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
