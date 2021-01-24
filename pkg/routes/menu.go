package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/uberswe/beubo/pkg/middleware"
	"github.com/uberswe/beubo/pkg/structs"
	"github.com/uberswe/beubo/pkg/utility"
	"net/http"
	"strconv"
)

// MenuAdmin shows menus that can be managed
func (br *BeuboRouter) MenuAdmin(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_menus", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// TODO implement menus
	params := mux.Vars(r)
	id := params["id"]

	i, err := strconv.Atoi(id)

	utility.ErrorHandler(err, false)

	site := structs.FetchSite(br.DB, i)

	self := r.Context().Value(middleware.UserContextKey)
	if self != nil && self.(structs.User).ID > 0 {
		if !self.(structs.User).CanAccessSite(br.DB, site) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}

	var pages []structs.Page

	extra := make(map[string]interface{})
	pagesRes := make(map[string]map[string]string)
	extra["SiteID"] = fmt.Sprintf("%d", site.ID)

	if err := br.DB.Where("site_id = ?", site.ID).Find(&pages).Error; err != nil {
		utility.ErrorHandler(err, false)
	}

	for _, page := range pages {
		pid := fmt.Sprintf("%d", page.ID)
		pagesRes[pid] = make(map[string]string)
		pagesRes[pid]["id"] = pid
		pagesRes[pid]["title"] = page.Title
		pagesRes[pid]["slug"] = page.Slug
	}
	extra["pagesRes"] = pagesRes

	pageData := structs.PageData{
		Template: "admin.site.page.home",
		Title:    "Admin",
		Extra:    extra,
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// AdminMenuAdd is the route for adding a menu
func (br *BeuboRouter) AdminMenuAdd(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_menus", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// TODO implement menus

	pageData := structs.PageData{
		Template: "admin.setting.add",
		Title:    "Admin - Add Setting",
		Themes:   br.Renderer.GetThemes(),
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// AdminMenuAddPost handles adding of a menu
func (br *BeuboRouter) AdminMenuAddPost(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_menus", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// TODO implement menus

	path := "/admin/settings/add"

	successMessage := "Setting created"
	invalidError := "an error occurred and the setting could not be created."

	key := r.FormValue("keyField")
	value := r.FormValue("valueField")

	if len(key) < 1 {
		invalidError = "The key is too short"
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, path, 302)
		return
	}
	if len(value) < 1 {
		invalidError = "The value is too short"
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, path, 302)
		return
	}

	if structs.CreateSetting(br.DB, key, value) {
		utility.SetFlash(w, "message", []byte(successMessage))
		http.Redirect(w, r, "/admin/settings", 302)
		return
	}

	utility.SetFlash(w, "error", []byte(invalidError))
	http.Redirect(w, r, "/admin/settings/add", 302)
}

// AdminMenuDelete handles the deletion of a menu
func (br *BeuboRouter) AdminMenuDelete(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_menus", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// TODO implement menus

	params := mux.Vars(r)
	id := params["id"]

	i, err := strconv.Atoi(id)

	utility.ErrorHandler(err, false)

	structs.DeleteSetting(br.DB, i)

	utility.SetFlash(w, "message", []byte("Setting deleted"))

	http.Redirect(w, r, "/admin/settings", 302)
}

// AdminMenuEdit is the route for editing a menu
func (br *BeuboRouter) AdminMenuEdit(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_menus", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// TODO implement menus

	params := mux.Vars(r)
	id := params["id"]

	i, err := strconv.Atoi(id)

	utility.ErrorHandler(err, false)

	setting := structs.FetchSetting(br.DB, i)

	if setting.ID == 0 {
		br.NotFoundHandler(w, r)
		return
	}

	pageData := structs.PageData{
		Template: "admin.setting.edit",
		Title:    "Admin - Edit Site",
		Extra:    setting,
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
	// TODO implement menus

	params := mux.Vars(r)
	id := params["id"]

	path := fmt.Sprintf("/admin/settings/edit/%s", id)

	i, err := strconv.Atoi(id)

	utility.ErrorHandler(err, false)

	successMessage := "Setting updated"
	invalidError := "an error occurred and the setting could not be updated."

	key := r.FormValue("keyField")
	value := r.FormValue("valueField")

	// TODO make rules for models
	if len(key) < 1 {
		invalidError = "The key is too short"
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, path, 302)
		return
	}
	if len(value) < 1 {
		invalidError = "The value is too short"
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, path, 302)
		return
	}

	if structs.UpdateSetting(br.DB, i, key, value) {
		utility.SetFlash(w, "message", []byte(successMessage))
		http.Redirect(w, r, "/admin/settings", 302)
		return
	}

	utility.SetFlash(w, "error", []byte(invalidError))
	http.Redirect(w, r, path, 302)
}
