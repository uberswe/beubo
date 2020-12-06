package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/uberswe/beubo/pkg/structs"
	"github.com/uberswe/beubo/pkg/utility"
	"net/http"
	"strconv"
)

// AdminUserAdd is the route for adding a site
func (br *BeuboRouter) AdminUserAdd(w http.ResponseWriter, r *http.Request) {
	pageData := structs.PageData{
		Template: "admin.user.add",
		Title:    "Admin - Add User",
		Themes:   br.Renderer.GetThemes(),
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// AdminUserAddPost handles adding of a global user
func (br *BeuboRouter) AdminUserAddPost(w http.ResponseWriter, r *http.Request) {
	path := "/admin/users/add"

	successMessage := "User created"
	invalidError := "an error occured and the user could not be created."

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

	if structs.CreateUser(br.DB, key, value) {
		utility.SetFlash(w, "message", []byte(successMessage))
		http.Redirect(w, r, "/admin/users", 302)
		return
	}

	utility.SetFlash(w, "error", []byte(invalidError))
	http.Redirect(w, r, "/admin/users/add", 302)
}

// AdminUserDelete handles the deletion of a global user
func (br *BeuboRouter) AdminUserDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	i, err := strconv.Atoi(id)

	utility.ErrorHandler(err, false)

	structs.DeleteUser(br.DB, i)

	utility.SetFlash(w, "message", []byte("User deleted"))

	http.Redirect(w, r, "/admin/users", 302)
}

// AdminUserEdit is the route for adding a user
func (br *BeuboRouter) AdminUserEdit(w http.ResponseWriter, r *http.Request) {
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

// AdminUserEditPost handles editing of a global user
func (br *BeuboRouter) AdminUserEditPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	path := fmt.Sprintf("/admin/users/edit/%s", id)

	i, err := strconv.Atoi(id)

	utility.ErrorHandler(err, false)

	successMessage := "User updated"
	invalidError := "an error occurred and the user could not be updated."

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

	if structs.UpdateUser(br.DB, i, key, value) {
		utility.SetFlash(w, "message", []byte(successMessage))
		http.Redirect(w, r, "/admin/users", 302)
		return
	}

	utility.SetFlash(w, "error", []byte(invalidError))
	http.Redirect(w, r, path, 302)
}
