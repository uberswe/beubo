package routes

import (
	"fmt"
	"github.com/markustenghamn/beubo/pkg/structs"
	"github.com/markustenghamn/beubo/pkg/utility"
	"net/http"
)

// Admin is the default admin route and template
func (br *BeuboRouter) Admin(w http.ResponseWriter, r *http.Request) {
	var sites []structs.Site

	extra := make(map[string]map[string]map[string]string)
	extra["sites"] = make(map[string]map[string]string)

	if err := br.DB.Find(&sites).Error; err != nil {
		utility.ErrorHandler(err, false)
	}

	for _, site := range sites {
		sid := fmt.Sprintf("%d", site.ID)
		extra["sites"][sid] = make(map[string]string)
		extra["sites"][sid]["id"] = sid
		extra["sites"][sid]["title"] = site.Title
		extra["sites"][sid]["domain"] = site.Domain
	}

	pageData := structs.PageData{
		Template: "admin.sites",
		Title:    "Admin",
		Extra:    extra,
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

func (br *BeuboRouter) Settings(w http.ResponseWriter, r *http.Request) {
	pageData := structs.PageData{
		Template: "admin.settings",
		Title:    "Settings",
	}
	br.Renderer.RenderHTMLPage(w, r, pageData)
}

func (br *BeuboRouter) Users(w http.ResponseWriter, r *http.Request) {
	pageData := structs.PageData{
		Template: "admin.users",
		Title:    "Users",
	}
	br.Renderer.RenderHTMLPage(w, r, pageData)

}

func (br *BeuboRouter) Plugins(w http.ResponseWriter, r *http.Request) {
	pageData := structs.PageData{
		Template: "admin.plugins",
		Title:    "Plugins",
	}
	br.Renderer.RenderHTMLPage(w, r, pageData)

}
