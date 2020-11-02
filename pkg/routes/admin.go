package routes

import (
	"fmt"
	"github.com/markustenghamn/beubo/pkg/structs"
	"github.com/markustenghamn/beubo/pkg/structs/page"
	"github.com/markustenghamn/beubo/pkg/structs/page/component"
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

	var rows []component.Row
	for _, site := range sites {
		sid := fmt.Sprintf("%d", site.ID)
		rows = append(rows, component.Row{
			Columns: []component.Column{
				{Name: "ID", Value: sid},
				{Name: "Site", Value: site.Title},
				{Name: "Domain", Value: site.Domain},
				{Name: ""},
				{Name: ""},
				{Name: ""},
				{Name: ""},
			},
		})
	}

	table := component.Table{
		Section: "main",
		Header: []component.Column{
			{Name: "ID"},
			{Name: "Site"},
			{Name: "Domain"},
			{Name: ""},
			{Name: ""},
			{Name: ""},
			{Name: ""},
		},
		Rows:             rows,
		PageNumber:       1,
		PageDisplayCount: 10,
		T:                br.Renderer.T,
	}

	pageData := structs.PageData{
		Template: "admin.sites",
		Title:    "Admin",
		Extra:    extra,
		Components: []page.Component{
			table,
		},
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
	var users []structs.User

	extra := make(map[string]map[string]map[string]string)
	extra["users"] = make(map[string]map[string]string)

	if err := br.DB.Find(&users).Error; err != nil {
		utility.ErrorHandler(err, false)
	}

	for _, user := range users {
		uid := fmt.Sprintf("%d", user.ID)
		extra["users"][uid] = make(map[string]string)
		extra["users"][uid]["id"] = uid
		extra["users"][uid]["email"] = user.Email
	}

	pageData := structs.PageData{
		Template: "admin.users",
		Title:    "Users",
		Extra:    extra,
	}
	br.Renderer.RenderHTMLPage(w, r, pageData)

}

func (br *BeuboRouter) GetPlugins(w http.ResponseWriter, r *http.Request) {

	extra := make(map[string]map[string]map[string]string)
	extra["plugins"] = *br.Plugins

	pageData := structs.PageData{
		Template: "admin.plugins",
		Title:    "Plugins",
		Extra:    extra,
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}
