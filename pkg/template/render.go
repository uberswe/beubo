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
}

func (btr *BeuboTemplateRenderer) Init() {
	log.Println("Parsing and loading templates...")
	funcMap := buildFuncMap()
	var err error
	btr.T, err = findAndParseTemplates(rootDir, funcMap)
	utility.ErrorHandler(err, false)
}

// RenderHTMLPage handles rendering of the html template and should be the last function called before returning the response
func (btr *BeuboTemplateRenderer) RenderHTMLPage(pageTitle string, pageTemplate string, w http.ResponseWriter, r *http.Request, extra interface{}) {

	if os.Getenv("ASSETS_DIR") != "" {
		rootDir = os.Getenv("ASSETS_DIR")
	}
	if btr.ReloadTemplates {
		log.Println("Parsing and loading templates...")
		funcMap := buildFuncMap()
		var err error
		btr.T, err = findAndParseTemplates(rootDir, funcMap)
		utility.ErrorHandler(err, false)

		if os.Getenv("THEME") != "" {
			btr.CurrentTheme = os.Getenv("THEME")
		}
	}

	var foundTemplate *template.Template

	path := fmt.Sprintf("%s.%s", currentTheme, pageTemplate)
	if foundTemplate = btr.T.Lookup(path); foundTemplate == nil {
		log.Printf("Theme file not found %s\n", path)
		return
	}

	// Session flash messages to prompt failed logins etc..
	errorMessage, err := utility.GetFlash(w, r, "error")
	utility.ErrorHandler(err, false)
	warningMessage, err := utility.GetFlash(w, r, "warning")
	utility.ErrorHandler(err, false)
	stringMessage, err := utility.GetFlash(w, r, "message")
	utility.ErrorHandler(err, false)

	data := structs.PageData{
		Title:       pageTitle,
		WebsiteName: "Beubo",
		URL:         "http://localhost:3000",
		Menu: []structs.MenuItem{
			{Title: "Home", Path: "/"},
			{Title: "Register", Path: "/register"},
			{Title: "Login", Path: "/login"},
		},
		Error:   string(errorMessage),
		Warning: string(warningMessage),
		Message: string(stringMessage),
		Year:    strconv.Itoa(time.Now().Year()),
		Extra:   extra,
	}

	if err != nil {
		log.Println("Could not serialize plugin message")
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
	log.Println("Reading from", cleanRoot)
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
