package beubo

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type TemplateRender struct {
	Template template.Template
	PageData PageData
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
