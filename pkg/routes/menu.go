package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/uberswe/beubo/pkg/middleware"
	"github.com/uberswe/beubo/pkg/structs"
	"github.com/uberswe/beubo/pkg/structs/page"
	"github.com/uberswe/beubo/pkg/structs/page/component"
	"github.com/uberswe/beubo/pkg/utility"
	"html/template"
	"net/http"
	"strconv"
)

// MenuAdmin shows menus that can be managed
func (br *BeuboRouter) MenuAdmin(w http.ResponseWriter, r *http.Request) {
	var err error
	if !middleware.CanAccess(br.DB, "manage_menus", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	// global menus have site_id 0
	i := 0
	if id != "" {
		i, err = strconv.Atoi(id)
		utility.ErrorHandler(err, false)

		if !currentUserCanAccessSite(id, br, r) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}

	var menus []structs.MenuSection

	extra := make(map[string]interface{})
	extra["SiteID"] = fmt.Sprintf("%d", i)

	if err := br.DB.Where("site_id = ?", i).Find(&menus).Error; err != nil {
		utility.ErrorHandler(err, false)
	}

	var rows []component.Row

	for _, menu := range menus {
		mid := fmt.Sprintf("%d", menu.ID)
		rows = append(rows, component.Row{
			Columns: []component.Column{
				{Name: "ID", Value: mid},
				{Name: "Section", Value: menu.Section},
				{Name: "", Field: component.Button{
					Link:    template.URL(fmt.Sprintf("/admin/menus/%s/edit", mid)),
					Class:   "btn btn-primary",
					Content: "Edit",
					T:       br.Renderer.T,
				}},
				{Name: "", Field: component.Button{
					Link:    template.URL(fmt.Sprintf("/admin/menus/%s/delete", mid)),
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
			{Name: "Section"},
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
		Title:    "Admin Menus",
		Components: []page.Component{
			component.Button{
				Section: "main",
				Link:    template.URL("/admin/menus/add"),
				Class:   "btn btn-primary",
				Content: "Add Menu",
				T:       br.Renderer.T,
			},
			table,
		},
		Extra: extra,
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// AdminMenuAdd is the route for adding a menu
func (br *BeuboRouter) AdminMenuAdd(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_menus", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	i, err := strconv.Atoi(id)
	utility.ErrorHandler(err, false)

	if !currentUserCanAccessSite(id, br, r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	extra := make(map[string]interface{})
	extra["SiteID"] = fmt.Sprintf("%d", i)

	pageData := structs.PageData{
		Template: "admin.menu.add",
		Title:    "Admin - Add Menu",
		Themes:   br.Renderer.GetThemes(),
		Extra:    extra,
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// AdminMenuAddPost handles adding of a menu
func (br *BeuboRouter) AdminMenuAddPost(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_menus", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	params := mux.Vars(r)
	id := params["id"]

	if !currentUserCanAccessSite(id, br, r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	path := "/admin/menus/add"

	successMessage := "Menu created"
	invalidError := "an error occurred and the menu could not be created."

	section := r.FormValue("sectionField")

	if len(section) < 1 {
		invalidError = "The section is too short"
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, path, 302)
		return
	}

	if menu := structs.CreateMenu(br.DB, section); menu.ID != 0 {
		utility.SetFlash(w, "message", []byte(successMessage))
		http.Redirect(w, r, "/admin/menus", 302)
		return
	}

	utility.SetFlash(w, "error", []byte(invalidError))
	http.Redirect(w, r, "/admin/sites/a/%s/menus/add", 302)
}

// AdminMenuDelete handles the deletion of a menu
func (br *BeuboRouter) AdminMenuDelete(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_menus", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	params := mux.Vars(r)
	menuId := params["menuId"]
	id := params["id"]
	i, err := strconv.Atoi(menuId)
	utility.ErrorHandler(err, false)

	if !currentUserCanAccessSite(id, br, r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	structs.DeleteMenu(br.DB, i)

	utility.SetFlash(w, "message", []byte("Menu deleted"))

	http.Redirect(w, r, "/admin/sites/a/%s/menus", 302)
}

// AdminMenuEdit is the route for editing a menu
func (br *BeuboRouter) AdminMenuEdit(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_menus", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	menuId := params["menuIid"]

	menuI, err := strconv.Atoi(menuId)
	utility.ErrorHandler(err, false)

	if !currentUserCanAccessSite(id, br, r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	menu := structs.FetchMenu(br.DB, menuI)

	if menu.ID == 0 {
		br.NotFoundHandler(w, r)
		return
	}

	pageData := structs.PageData{
		Template: "admin.menu.edit",
		Title:    "Admin - Edit Menu",
		Extra:    menu,
		Themes:   br.Renderer.GetThemes(),
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// AdminMenuEditPost handles editing of a menu
func (br *BeuboRouter) AdminMenuEditPost(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_menus", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	params := mux.Vars(r)
	id := params["id"]

	path := fmt.Sprintf("/admin/sites/a/%s/menus/edit/%s", id)

	i, err := strconv.Atoi(id)

	utility.ErrorHandler(err, false)

	successMessage := "Menu updated"
	invalidError := "an error occurred and the menu could not be updated."

	section := r.FormValue("sectionField")

	// TODO make rules for models
	if len(section) < 1 {
		invalidError = "The key is too short"
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, path, 302)
		return
	}

	if structs.UpdateMenu(br.DB, i, section) {
		utility.SetFlash(w, "message", []byte(successMessage))
		http.Redirect(w, r, "/admin/sites/a/%s/menus", 302)
		return
	}

	utility.SetFlash(w, "error", []byte(invalidError))
	http.Redirect(w, r, path, 302)
}
