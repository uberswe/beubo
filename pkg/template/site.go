package template

import (
	"io/ioutil"
	"log"
)

// GetThemes fetches all the themes in the theme directory as strings
func (btr *BeuboTemplateRenderer) GetThemes() map[string]string {
	pageThemes := map[string]string{}

	files, err := ioutil.ReadDir(btr.ThemeDir)

	if err != nil {
		log.Println(err)
		return pageThemes
	}

	for _, file := range files {
		// TODO in the future we could define a way to define names for themes
		// Ignore the install directory, only used for installation
		if file.IsDir() && file.Name() != "install" {
			pageThemes[file.Name()] = file.Name()
		}
	}

	return pageThemes
}
