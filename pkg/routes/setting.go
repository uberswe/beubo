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

// AdminSettingAdd is the route for adding a site
func (br *BeuboRouter) AdminSettingAdd(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_settings", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	pageData := structs.PageData{
		Template: "admin.setting.add",
		Title:    "Admin - Add Setting",
		Themes:   br.Renderer.GetThemes(),
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// AdminSettingAddPost handles adding of a global setting
func (br *BeuboRouter) AdminSettingAddPost(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_settings", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

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

// AdminSettingDelete handles the deletion of a global setting
func (br *BeuboRouter) AdminSettingDelete(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_settings", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	params := mux.Vars(r)
	id := params["id"]

	i, err := strconv.Atoi(id)

	utility.ErrorHandler(err, false)

	structs.DeleteSetting(br.DB, i)

	utility.SetFlash(w, "message", []byte("Setting deleted"))

	http.Redirect(w, r, "/admin/settings", 302)
}

// AdminSettingEdit is the route for adding a setting
func (br *BeuboRouter) AdminSettingEdit(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_settings", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

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

// AdminSettingEditPost handles editing of a global setting
func (br *BeuboRouter) AdminSettingEditPost(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_settings", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

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
