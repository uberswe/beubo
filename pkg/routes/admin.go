package routes

import (
	"fmt"
	"github.com/uberswe/beubo/pkg/middleware"
	"github.com/uberswe/beubo/pkg/plugin"
	"github.com/uberswe/beubo/pkg/structs"
	"github.com/uberswe/beubo/pkg/structs/page"
	"github.com/uberswe/beubo/pkg/structs/page/component"
	"github.com/uberswe/beubo/pkg/utility"
	"html/template"
	"net/http"
)

// Admin is the default admin route and template
func (br *BeuboRouter) Admin(w http.ResponseWriter, r *http.Request) {
	var sites []structs.Site

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
				{Name: "", Field: component.Button{
					// TODO fix schema here
					Link:    template.URL(fmt.Sprintf("%s://%s/", "http", site.Domain)),
					Class:   "btn btn-primary",
					Content: "View",
					T:       br.Renderer.T,
				}},
				{Name: "", Field: component.Button{
					Link:    template.URL(fmt.Sprintf("/admin/sites/a/%s", sid)),
					Class:   "btn btn-primary",
					Content: "Manage",
					T:       br.Renderer.T,
				}},
				{Name: "", Field: component.Button{
					Link:    template.URL(fmt.Sprintf("/admin/sites/edit/%s", sid)),
					Class:   "btn btn-primary",
					Content: "Edit",
					T:       br.Renderer.T,
				}},
				{Name: "", Field: component.Button{
					Link:    template.URL(fmt.Sprintf("/admin/sites/delete/%s", sid)),
					Class:   "btn btn-primary",
					Content: "Delete",
					T:       br.Renderer.T,
				}},
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
		Template: "admin.page",
		Title:    "Admin - Sites",
		Components: []page.Component{
			component.Button{
				Section: "main",
				Link:    template.URL("/admin/sites/add"),
				Class:   "btn btn-primary",
				Content: "Add Site",
				T:       br.Renderer.T,
			},
			table,
		},
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// Settings is the route for loading the admin settings page
func (br *BeuboRouter) Settings(w http.ResponseWriter, r *http.Request) {
	var settings []structs.Setting

	if err := br.DB.Find(&settings).Error; err != nil {
		utility.ErrorHandler(err, false)
	}

	var rows []component.Row
	for _, setting := range settings {
		sid := fmt.Sprintf("%d", setting.ID)
		rows = append(rows, component.Row{
			Columns: []component.Column{
				{Name: "ID", Value: sid},
				{Name: "Site", Value: setting.Key},
				{Name: "Domain", Value: setting.Value},
				{Name: "", Field: component.Button{
					Link:    template.URL(fmt.Sprintf("/admin/settings/edit/%s", sid)),
					Class:   "btn btn-primary",
					Content: "Edit",
					T:       br.Renderer.T,
				}},
				{Name: "", Field: component.Button{
					Link:    template.URL(fmt.Sprintf("/admin/settings/delete/%s", sid)),
					Class:   "btn btn-primary",
					Content: "Delete",
					T:       br.Renderer.T,
				}},
			},
		})
	}

	table := component.Table{
		Section: "main",
		Header: []component.Column{
			{Name: "ID"},
			{Name: "Key"},
			{Name: "Value"},
			{Name: ""},
			{Name: ""},
		},
		Rows:             rows,
		PageNumber:       1,
		PageDisplayCount: 10,
		T:                br.Renderer.T,
	}

	pageData := structs.PageData{
		Template: "admin.page",
		Title:    "Admin - Settings",
		Components: []page.Component{
			component.Button{
				Section: "main",
				Link:    template.URL("/admin/settings/add"),
				Class:   "btn btn-primary",
				Content: "Add Setting",
				T:       br.Renderer.T,
			},
			table,
			component.Text{
				Section: "main",
				Content: "This is the very early stages of how settings works. Currently it is just plain key and values with no checks.<br/>" +
					"To use IP whitelisting first set any IPs to whitelist with <b>whitelisted_ip</b> as the key and the IP as value.<br/>" +
					"Then set <b>ip_whitelist</b> to <b>true</b> to enable the whitelist blocking",
				T: br.Renderer.T,
			},
		},
	}
	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// Users is the route for loading the admin users page
func (br *BeuboRouter) Users(w http.ResponseWriter, r *http.Request) {
	var users []structs.User

	if err := br.DB.Find(&users).Error; err != nil {
		utility.ErrorHandler(err, false)
	}

	self := r.Context().Value(middleware.UserContextKey)

	var rows []component.Row
	if self != nil && self.(structs.User).ID > 0 {
		for _, user := range users {
			sid := fmt.Sprintf("%d", user.ID)
			if self.(structs.User).ID == user.ID {
				rows = append(rows, component.Row{
					Columns: []component.Column{
						{Name: "ID", Value: sid},
						{Name: "Email", Value: user.Email},
						{},
						{},
					},
				})
			} else {
				rows = append(rows, component.Row{
					Columns: []component.Column{
						{Name: "ID", Value: sid},
						{Name: "Email", Value: user.Email},
						{Name: "", Field: component.Button{
							Link:    template.URL(fmt.Sprintf("/admin/users/edit/%s", sid)),
							Class:   "btn btn-primary",
							Content: "Edit",
							T:       br.Renderer.T,
						}},
						{Name: "", Field: component.Button{
							Link:    template.URL(fmt.Sprintf("/admin/users/delete/%s", sid)),
							Class:   "btn btn-primary",
							Content: "Delete",
							T:       br.Renderer.T,
						}},
					},
				})
			}
		}
	}

	table := component.Table{
		Section: "main",
		Header: []component.Column{
			{Name: "ID"},
			{Name: "Email"},
			{Name: ""},
			{Name: ""},
		},
		Rows:             rows,
		PageNumber:       1,
		PageDisplayCount: 10,
		T:                br.Renderer.T,
	}

	pageData := structs.PageData{
		Template: "admin.page",
		Title:    "Admin - Users",
		Components: []page.Component{
			component.Button{
				Section: "main",
				Link:    template.URL("/admin/users/add"),
				Class:   "btn btn-primary",
				Content: "Add User",
				T:       br.Renderer.T,
			},
			table,
		},
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)

}

// Plugins is the route for loading the admin plugins page
func (br *BeuboRouter) Plugins(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// TODO handle post
		return
	}

	// reload plugins when viewing the plugin page
	plh := plugin.Load(*br.PluginHandler)
	br.PluginHandler = &plh

	var rows []component.Row
	for _, p := range br.PluginHandler.Plugins {
		comprow := component.Row{
			Columns: []component.Column{
				{Name: "Name", Value: fmt.Sprintf("<a href=\"/admin/plugins/edit/%s\">%s</a>", p.Definition, p.Definition)},
			},
		}
		rows = append(rows, comprow)
	}

	table := component.Table{
		Section: "main",
		Header: []component.Column{
			{Name: "Name"},
		},
		Rows:             rows,
		PageNumber:       1,
		PageDisplayCount: 10,
		T:                br.Renderer.T,
	}

	pageData := structs.PageData{
		Template: "admin.page",
		Title:    "Admin - Plugins",
		Components: []page.Component{
			table,
		},
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}
