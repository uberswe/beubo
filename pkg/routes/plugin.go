package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/uberswe/beubo/pkg/middleware"
	"github.com/uberswe/beubo/pkg/plugin"
	"github.com/uberswe/beubo/pkg/structs"
	"github.com/uberswe/beubo/pkg/structs/page"
	"github.com/uberswe/beubo/pkg/structs/page/component"
	"github.com/uberswe/beubo/pkg/utility"
	"html/template"
	"net/http"
)

// AdminPluginEdit is the route for editing a plugin
func (br *BeuboRouter) AdminPluginEdit(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_plugins", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	params := mux.Vars(r)
	id := params["id"]

	plugins := plugin.FetchPluginSites(br.DB, id)

	if len(plugins) <= 0 {
		br.NotFoundHandler(w, r)
		return
	}

	var rows []component.Row
	for _, p := range plugins {
		comprow := component.Row{
			Columns: []component.Column{
				{Name: "Name", Value: p.Site.Title},
				{Field: component.CheckBoxField{
					Name:       fmt.Sprintf("%s[]", p.PluginIdentifier),
					Identifier: fmt.Sprintf("%s_%d", p.PluginIdentifier, p.Site.ID),
					Value:      fmt.Sprintf("%d", p.SiteID),
					Checked:    p.Active,
					T:          br.Renderer.T,
				}},
			},
		}
		rows = append(rows, comprow)
	}

	button := component.Button{
		Section: "main",
		Link:    template.URL("/admin/plugins"),
		Class:   "btn btn-primary",
		Content: "Back",
		T:       br.Renderer.T,
	}

	formButton := component.Button{
		Section: "main",
		Class:   "btn btn-primary",
		Content: "Save",
		T:       br.Renderer.T,
	}

	table := component.Table{
		Section: "main",
		Header: []component.Column{
			{Name: "Site"},
			{Name: "Active"},
		},
		Rows:             rows,
		PageNumber:       1,
		PageDisplayCount: 10,
		T:                br.Renderer.T,
	}

	form := component.Form{
		Section: "main",
		Fields: []page.Component{
			table,
			formButton,
		},
		T:      br.Renderer.T,
		Method: "POST",
		Action: fmt.Sprintf("/admin/plugins/edit/%s", id),
	}

	pageData := structs.PageData{
		Template: "admin.page",
		Title:    "Admin - Edit Plugin",
		Components: []page.Component{
			button,
			form,
		},
		Themes: br.Renderer.GetThemes(),
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// AdminPluginEditPost handles editing of plugins
func (br *BeuboRouter) AdminPluginEditPost(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_plugins", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err := r.ParseForm()
	utility.ErrorHandler(err, false)

	params := mux.Vars(r)
	id := params["id"]

	plugins := plugin.FetchPluginSites(br.DB, id)

	path := fmt.Sprintf("/admin/plugins/edit/%s", id)

	successMessage := "Plugin updated"

	for key, values := range r.PostForm {
		if key == fmt.Sprintf("%s[]", id) {
			for _, p := range plugins {
				found := false
				for _, v := range values {
					if v == fmt.Sprintf("%d", p.SiteID) {
						found = true
					}
				}
				if found {
					// The plugin is active
					p.Active = true
				} else {
					// The plugin is inactive
					p.Active = false
				}
				br.DB.Save(&p)
			}
		}
	}

	// If all sites are set to inactive then we don't get any post data
	if len(r.PostForm) <= 0 {
		for _, p := range plugins {
			p.Active = false
			br.DB.Save(&p)
		}
	}

	utility.SetFlash(w, "message", []byte(successMessage))
	http.Redirect(w, r, path, 302)
}
