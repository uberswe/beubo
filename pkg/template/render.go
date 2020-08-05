package template

import (
	"fmt"
	"github.com/markustenghamn/beubo/pkg/structs"
	"github.com/markustenghamn/beubo/pkg/utility"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	currentTheme = "default"
	rootDir      = "./web/"
)

type BeuboTemplateRenderer struct {
	T               *template.Template
	ReloadTemplates bool
	CurrentTheme    string
	ThemeDir        string
}

func (btr *BeuboTemplateRenderer) Init() {
	log.Println("Parsing and loading templates...")
	funcMap := buildFuncMap()
	var err error
	btr.T, err = findAndParseTemplates(rootDir, funcMap)
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
	site := r.Context().Value("site")
	user := r.Context().Value("user")
	if site != nil {
		btr.CurrentTheme = site.(structs.Site).Theme.Slug
		siteName = site.(structs.Site).Title
	} else if os.Getenv("THEME") != "" {
		// Default theme
		btr.CurrentTheme = os.Getenv("THEME")
	} else {
		// Default theme
		btr.CurrentTheme = "default"
	}

	themePath := path.Join(rootDir, "/themes/")

	var stylesheets []string

	err = filepath.Walk(path.Join(themePath, btr.CurrentTheme, "/css/"), func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		p = strings.TrimLeft(p, themePath)
		li := strings.LastIndex(p, "/")
		if strings.HasSuffix(p, ".css") && strings.Contains(r.URL.Path, strings.Replace(p[:li], path.Join(btr.CurrentTheme, "/css/"), "", 1)) {
			stylesheets = append(stylesheets, "/"+p)
		}
		return nil
	})

	if err != nil {
		log.Println(err)
	}

	var sToDelete []int
	for si, s := range stylesheets {
		if !strings.Contains(s, ".min.css") {
			for si2, s2 := range stylesheets {
				if si != si2 && strings.TrimSuffix(s, ".css") == strings.TrimSuffix(s2, ".min.css") {
					sToDelete = append(sToDelete, si)
					break
				}
			}
		}
	}

	removed := 0
	for _, si := range sToDelete {
		if len(stylesheets) > (si - removed) {
			stylesheets = append(stylesheets[:(si-removed)], stylesheets[(si-removed+1):]...)
		} else {
			stylesheets = stylesheets[:(si - removed)]
		}
		removed++
	}

	menu := []structs.MenuItem{
		{Title: "Home", Path: "/"},
		{Title: "Login", Path: "/login"},
		{Title: "Register", Path: "/register"},
	}

	if user != nil && user.(structs.User).ID > 0 {
		menu = []structs.MenuItem{
			{Title: "Home", Path: "/"},
			{Title: "Admin", Path: "/admin"},
			{Title: "Logout", Path: "/logout"},
		}
	}

	data := structs.PageData{
		Stylesheets: stylesheets,
		// TODO make the favicon dynamic
		Favicon:     "/default/images/favicon.ico",
		WebsiteName: siteName,
		URL:         "http://localhost:3000",
		// TODO make the menu dynamic
		Menu:    menu,
		Error:   string(errorMessage),
		Warning: string(warningMessage),
		Message: string(stringMessage),
		Year:    strconv.Itoa(time.Now().Year()),
	}

	data = mergePageData(data, pageData)

	if btr.ReloadTemplates {
		log.Println("Parsing and loading templates...")
		funcMap := buildFuncMap()
		var err error
		btr.T, err = findAndParseTemplates(rootDir, funcMap)
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
func findAndParseTemplates(rootDir string, funcMap template.FuncMap) (*template.Template, error) {
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
			t := root.New(name).Funcs(funcMap)
			t, e2 = t.Parse(string(b))
			if e2 != nil {
				return e2
			}
		}

		return nil
	})

	return root, err
}

// TODO we could build function maps for some areas of the template? but why?
func buildFuncMap() template.FuncMap {
	return template.FuncMap{
		"bContent": func(feature string) bool {
			return false
		},
		"bHead": func(feature string) bool {
			return false
		},
		"bHeader": func(feature string) bool {
			return false
		},
		"bFooter": func(feature string) bool {
			return false
		},
	}
}

func mergePageData(a structs.PageData, b structs.PageData) structs.PageData {
	// TODO this could be simplified by making a function that compares an interface and picks a value but I decided that this is more readable for now
	if b.Template != "" {
		a.Template = b.Template
	}

	if b.Title != "" {
		a.Title = b.Title
	}

	if b.Content != "" {
		a.Content = b.Content
	}

	if b.WebsiteName != "" {
		a.WebsiteName = b.WebsiteName
	}

	if b.Menu != nil {
		a.Menu = b.Menu
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

	a.Themes = b.Themes
	a.Templates = b.Templates
	a.Extra = b.Extra

	return a
}
