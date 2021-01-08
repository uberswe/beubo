package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/uberswe/beubo/pkg/structs"
	"github.com/uberswe/beubo/pkg/structs/page"
	"github.com/uberswe/beubo/pkg/structs/page/component"
	"github.com/uberswe/beubo/pkg/utility"
	"html/template"
	"net/http"
	"strconv"
)

// AdminUserAdd is the route for adding a user
func (br *BeuboRouter) AdminUserAdd(w http.ResponseWriter, r *http.Request) {
	pageData := structs.PageData{
		Template: "admin.user.add",
		Title:    "Admin - Add User",
		Themes:   br.Renderer.GetThemes(),
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// AdminUserAddPost handles adding of a user
func (br *BeuboRouter) AdminUserAddPost(w http.ResponseWriter, r *http.Request) {
	path := "/admin/users/add"

	successMessage := "User created"
	invalidError := "an error occurred and the user could not be created."

	email := r.FormValue("emailField")
	password := r.FormValue("passwordField")

	if !utility.IsEmailValid(email) {
		invalidError = "The email is invalid"
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, path, 302)
		return
	}
	if len(password) < 8 {
		invalidError = "The password is too short"
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, path, 302)
		return
	}

	if structs.CreateUser(br.DB, email, password) {
		utility.SetFlash(w, "message", []byte(successMessage))
		http.Redirect(w, r, "/admin/users", 302)
		return
	}

	utility.SetFlash(w, "error", []byte(invalidError))
	http.Redirect(w, r, "/admin/users/add", 302)
}

// AdminUserDelete handles the deletion of a user
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

// AdminUserEditPost handles editing of a user
func (br *BeuboRouter) AdminUserEditPost(w http.ResponseWriter, r *http.Request) {
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

// AdminUserRoles is the route for managing user roles
func (br *BeuboRouter) AdminUserRoles(w http.ResponseWriter, r *http.Request) {
	var roles []structs.Role

	if err := br.DB.Find(&roles).Error; err != nil {
		utility.ErrorHandler(err, false)
	}

	var rows []component.Row
	for _, role := range roles {
		sid := fmt.Sprintf("%d", role.ID)
		// TODO For now it's not possible to edit default roles. This is because they would be added again when launching Beubo and this needs to be handled in a better way.
		if role.Name == "Administrator" || role.Name == "Member" {
			rows = append(rows, component.Row{
				Columns: []component.Column{
					{Name: "ID", Value: sid},
					{Name: "Name", Value: role.Name},
					{},
					{},
				},
			})
		} else {
			rows = append(rows, component.Row{
				Columns: []component.Column{
					{Name: "ID", Value: sid},
					{Name: "Name", Value: role.Name},
					{Name: "", Field: component.Button{
						Link:    template.URL(fmt.Sprintf("/admin/users/roles/edit/%s", sid)),
						Class:   "btn btn-primary",
						Content: "Edit",
						T:       br.Renderer.T,
					}},
					{Name: "", Field: component.Button{
						Link:    template.URL(fmt.Sprintf("/admin/users/roles/delete/%s", sid)),
						Class:   "btn btn-primary",
						Content: "Delete",
						T:       br.Renderer.T,
					}},
				},
			})
		}
	}

	table := component.Table{
		Section: "main",
		Header: []component.Column{
			{Name: "ID"},
			{Name: "Name"},
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
		Title:    "Admin - User Roles",
		Components: []page.Component{
			component.Button{
				Section: "main",
				Link:    template.URL("/admin/users/roles/add"),
				Class:   "btn btn-primary",
				Content: "Add Role",
				T:       br.Renderer.T,
			},
			table,
		},
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// AdminUserRoleAdd is the route for adding a user role
func (br *BeuboRouter) AdminUserRoleAdd(w http.ResponseWriter, r *http.Request) {
	pageData := structs.PageData{
		Template: "admin.user.role.add",
		Title:    "Admin - Add Role",
		Themes:   br.Renderer.GetThemes(),
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// AdminUserRoleAddPost handles adding of a user role
func (br *BeuboRouter) AdminUserRoleAddPost(w http.ResponseWriter, r *http.Request) {
	path := "/admin/users/roles/add"

	successMessage := "Role created"
	invalidError := "an error occurred and the user could not be created."

	name := r.FormValue("nameField")

	if len(name) < 2 {
		invalidError = "The name is too short"
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, path, 302)
		return
	}

	if structs.CreateRole(br.DB, name) {
		utility.SetFlash(w, "message", []byte(successMessage))
		http.Redirect(w, r, "/admin/users/roles", 302)
		return
	}

	utility.SetFlash(w, "error", []byte(invalidError))
	http.Redirect(w, r, "/admin/users/roles/add", 302)
}

// AdminUserRoleDelete handles the deletion of a user role
func (br *BeuboRouter) AdminUserRoleDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	i, err := strconv.Atoi(id)

	utility.ErrorHandler(err, false)

	role := structs.FetchRole(br.DB, i)
	if role.ID != 0 && !role.IsDefault() {
		structs.DeleteRole(br.DB, i)
	} else {
		utility.SetFlash(w, "error", []byte("Can not delete a default role"))
		http.Redirect(w, r, "/admin/users/roles", 302)
		return
	}

	utility.SetFlash(w, "message", []byte("Role deleted"))

	http.Redirect(w, r, "/admin/users/roles", 302)
}

// AdminUserRoleEdit is the route for adding a user role
func (br *BeuboRouter) AdminUserRoleEdit(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	i, err := strconv.Atoi(id)

	utility.ErrorHandler(err, false)

	role := structs.FetchRole(br.DB, i)
	if role.ID == 0 || role.IsDefault() {
		br.NotFoundHandler(w, r)
		return
	}

	pageData := structs.PageData{
		Template: "admin.user.role.edit",
		Title:    "Admin - Edit Role",
		Extra:    role,
		Themes:   br.Renderer.GetThemes(),
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// AdminUserRoleEditPost handles editing of a user role
func (br *BeuboRouter) AdminUserRoleEditPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	path := fmt.Sprintf("/admin/users/roles/edit/%s", id)

	i, err := strconv.Atoi(id)

	utility.ErrorHandler(err, false)

	successMessage := "Role updated"
	invalidError := "an error occurred and the role could not be updated."

	name := r.FormValue("nameField")

	if len(name) < 2 {
		invalidError = "The name is too short"
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, path, 302)
		return
	}

	if structs.UpdateRole(br.DB, i, name) {
		utility.SetFlash(w, "message", []byte(successMessage))
		http.Redirect(w, r, "/admin/users/roles", 302)
		return
	}

	utility.SetFlash(w, "error", []byte(invalidError))
	http.Redirect(w, r, path, 302)
}
