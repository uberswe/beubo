package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/uberswe/beubo/pkg/structs"
	"github.com/uberswe/beubo/pkg/utility"
	"net/http"
	"strconv"
)

// AdminPluginEdit is the route for editing a plugin
func (br *BeuboRouter) AdminPluginEdit(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	i, err := strconv.Atoi(id)

	utility.ErrorHandler(err, false)

	user := structs.FetchUser(br.DB, i)

	if user.ID == 0 {
		br.NotFoundHandler(w, r)
		return
	}

	pageData := structs.PageData{
		Template: "admin.user.edit",
		Title:    "Admin - Edit Site",
		Extra:    user,
		Themes:   br.Renderer.GetThemes(),
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// AdminPluginEditPost handles editing of plugins
func (br *BeuboRouter) AdminPluginEditPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	path := fmt.Sprintf("/admin/users/edit/%s", id)

	i, err := strconv.Atoi(id)

	utility.ErrorHandler(err, false)

	successMessage := "User updated"
	invalidError := "an error occurred and the user could not be updated."

	// TODO make rules for models
	email := r.FormValue("emailField")
	password := r.FormValue("passwordField")

	if !utility.IsEmailValid(email) {
		invalidError = "The email is invalid"
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, path, 302)
		return
	}
	if len(password) > 0 && len(password) < 8 {
		invalidError = "The password is too short"
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, path, 302)
		return
	}

	if structs.UpdateUser(br.DB, i, email, password) {
		utility.SetFlash(w, "message", []byte(successMessage))
		http.Redirect(w, r, "/admin/users", 302)
		return
	}

	utility.SetFlash(w, "error", []byte(invalidError))
	http.Redirect(w, r, path, 302)
}
