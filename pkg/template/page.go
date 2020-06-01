package template

import (
	"strings"
)

func (btr *BeuboTemplateRenderer) GetPageTemplates() map[string]string {
	pageTemplates := map[string]string{
		"page": "Default page",
	}
	for _, t := range btr.T.Templates() {
		// Ignore paths and only look at names containing .page. and does not contain .admin.
		if !strings.Contains(t.Name(), "themes/") && !strings.Contains(t.Name(), ".admin.") && strings.Contains(t.Name(), ".page.") {
			s := strings.Split(t.Name(), ".")
			pageName := ""
			pageTemplate := "page"
			pageFound := false
			for _, part := range s {
				if pageFound {
					if len(pageName) > 0 {
						pageName += " "
					}
					pageTemplate += "." + part
					pageName += strings.Title(part)
				}
				// Anything after page is our page template name
				if part == "page" {
					pageFound = true
				}
			}
			if len(pageName) > 0 {
				pageTemplates[pageTemplate] = pageName
			}
		}
	}
	return pageTemplates
}
