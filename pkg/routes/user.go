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

type userAddData struct {
	Roles []structs.Role
	Sites []structs.Site
}

// AdminUserAdd is the route for adding a user
func (br *BeuboRouter) AdminUserAdd(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_users", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	roles := []structs.Role{}
	sites := []structs.Site{}
	br.DB.Find(&roles)
	br.DB.Find(&sites)
	pageData := structs.PageData{
		Template: "admin.user.add",
		Title:    "Admin - Add User",
		Themes:   br.Renderer.GetThemes(),
		Extra: userAddData{
			Roles: roles,
			Sites: sites,
		},
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// AdminUserAddPost handles adding of a user
func (br *BeuboRouter) AdminUserAddPost(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_users", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	path := "/admin/users/add"

	successMessage := "User created"
	invalidError := "an error occurred and the user could not be created."

	email := r.FormValue("emailField")
	password := r.FormValue("passwordField")

	rs := []int{}
	ss := []int{}

	for key, values := range r.PostForm {
		if key == fmt.Sprintf("%s[]", "roleField") {
			for _, value := range values {
				valueInt, err := strconv.Atoi(value)
				utility.ErrorHandler(err, false)
				rs = append(rs, valueInt)
			}
		} else if key == fmt.Sprintf("%s[]", "siteField") {
			for _, value := range values {
				valueInt, err := strconv.Atoi(value)
				utility.ErrorHandler(err, false)
				ss = append(ss, valueInt)
			}
		}
	}

	roles := []*structs.Role{}
	if len(rs) > 0 {
		br.DB.Find(&roles, rs)
	}
	sites := []*structs.Site{}
	if len(ss) > 0 {
		br.DB.Find(&sites, ss)
	}

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

	if structs.CreateUser(br.DB, email, password, roles, sites) {
		utility.SetFlash(w, "message", []byte(successMessage))
		http.Redirect(w, r, "/admin/users", 302)
		return
	}

	utility.SetFlash(w, "error", []byte(invalidError))
	http.Redirect(w, r, "/admin/users/add", 302)
}

// AdminUserDelete handles the deletion of a user
func (br *BeuboRouter) AdminUserDelete(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_users", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	params := mux.Vars(r)
	id := params["id"]

	i, err := strconv.Atoi(id)

	utility.ErrorHandler(err, false)

	structs.DeleteUser(br.DB, i)

	utility.SetFlash(w, "message", []byte("User deleted"))

	http.Redirect(w, r, "/admin/users", 302)
}

type userEdit struct {
	User     structs.User
	Role     structs.Role
	Roles    []userEditRole
	Features []userEditFeature
	Sites    []userEditSite
}

type userEditRole struct {
	Role    structs.Role
	Checked bool
}

type userEditSite struct {
	Site    structs.Site
	Checked bool
}

// AdminUserEdit is the route for adding a user
func (br *BeuboRouter) AdminUserEdit(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_users", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	params := mux.Vars(r)
	id := params["id"]

	i, err := strconv.Atoi(id)

	utility.ErrorHandler(err, false)

	user := structs.FetchUser(br.DB, i)

	if user.ID == 0 {
		br.NotFoundHandler(w, r)
		return
	}

	ueRoles := []userEditRole{}
	roles := []structs.Role{}
	br.DB.Find(&roles)
	for _, role := range roles {
		ueRoles = append(ueRoles, userEditRole{Role: role, Checked: user.HasRole(br.DB, role)})
	}

	ueSites := []userEditSite{}
	sites := []structs.Site{}
	br.DB.Find(&sites)
	for _, site := range sites {
		ueSites = append(ueSites, userEditSite{Site: site, Checked: user.CanAccessSite(br.DB, site)})
	}

	pageData := structs.PageData{
		Template: "admin.user.edit",
		Title:    "Admin - Edit Site",
		Extra: userEdit{
			User:  user,
			Roles: ueRoles,
			Sites: ueSites,
		},
		Themes: br.Renderer.GetThemes(),
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// AdminUserEditPost handles editing of a user
func (br *BeuboRouter) AdminUserEditPost(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_users", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

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

	rs := []int{}
	ss := []int{}

	for key, values := range r.PostForm {
		if key == fmt.Sprintf("%s[]", "roleField") {
			for _, value := range values {
				valueInt, err := strconv.Atoi(value)
				utility.ErrorHandler(err, false)
				rs = append(rs, valueInt)
			}
		} else if key == fmt.Sprintf("%s[]", "siteField") {
			for _, value := range values {
				valueInt, err := strconv.Atoi(value)
				utility.ErrorHandler(err, false)
				ss = append(ss, valueInt)
			}
		}
	}

	roles := []*structs.Role{}
	if len(rs) > 0 {
		br.DB.Find(&roles, rs)
	}
	sites := []*structs.Site{}
	if len(ss) > 0 {
		br.DB.Find(&sites, ss)
	}

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

	if structs.UpdateUser(br.DB, i, email, password, roles, sites) {
		utility.SetFlash(w, "message", []byte(successMessage))
		http.Redirect(w, r, "/admin/users", 302)
		return
	}

	utility.SetFlash(w, "error", []byte(invalidError))
	http.Redirect(w, r, path, 302)
}

// AdminUserRoles is the route for managing user roles
func (br *BeuboRouter) AdminUserRoles(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_users", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

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
					{Name: "Name", Value: role.Name},
					{},
					{},
				},
			})
		} else {
			rows = append(rows, component.Row{
				Columns: []component.Column{
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
	if !middleware.CanAccess(br.DB, "manage_user_roles", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	features := []structs.Feature{}
	br.DB.Find(&features)

	pageData := structs.PageData{
		Template: "admin.user.role.add",
		Title:    "Admin - Add Role",
		Themes:   br.Renderer.GetThemes(),
		Extra:    features,
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// AdminUserRoleAddPost handles adding of a user role
func (br *BeuboRouter) AdminUserRoleAddPost(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_user_roles", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

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
	if !middleware.CanAccess(br.DB, "manage_user_roles", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

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

type userEditFeature struct {
	Feature structs.Feature
	Checked bool
}

// AdminUserRoleEdit is the route for adding a user role
func (br *BeuboRouter) AdminUserRoleEdit(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_user_roles", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	params := mux.Vars(r)
	id := params["id"]

	i, err := strconv.Atoi(id)

	utility.ErrorHandler(err, false)

	role := structs.FetchRole(br.DB, i)
	if role.ID == 0 || role.IsDefault() {
		br.NotFoundHandler(w, r)
		return
	}

	ueFeatures := []userEditFeature{}
	features := []structs.Feature{}
	br.DB.Find(&features)
	for _, feature := range features {
		ueFeatures = append(ueFeatures, userEditFeature{Feature: feature, Checked: role.HasFeature(br.DB, feature)})
	}

	pageData := structs.PageData{
		Template: "admin.user.role.edit",
		Title:    "Admin - Edit Role",
		Extra:    userEdit{Features: ueFeatures, Role: role},
		Themes:   br.Renderer.GetThemes(),
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// AdminUserRoleEditPost handles editing of a user role
func (br *BeuboRouter) AdminUserRoleEditPost(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_user_roles", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err := r.ParseForm()
	utility.ErrorHandler(err, false)
	params := mux.Vars(r)
	id := params["id"]

	path := fmt.Sprintf("/admin/users/roles/edit/%s", id)

	i, err := strconv.Atoi(id)

	utility.ErrorHandler(err, false)

	successMessage := "Role updated"
	invalidError := "an error occurred and the role could not be updated."

	fs := []int{}

	name := r.FormValue("nameField")
	for key, values := range r.PostForm {
		if key == fmt.Sprintf("%s[]", "featureField") {
			for _, value := range values {
				valueInt, err := strconv.Atoi(value)
				utility.ErrorHandler(err, false)
				fs = append(fs, valueInt)
			}
		}
	}

	features := []*structs.Feature{}
	if len(fs) > 0 {
		br.DB.Find(&features, fs)
	}

	if len(name) < 2 {
		invalidError = "The name is too short"
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, path, 302)
		return
	}

	if structs.UpdateRole(br.DB, i, name, features) {
		utility.SetFlash(w, "message", []byte(successMessage))
		http.Redirect(w, r, "/admin/users/roles", 302)
		return
	}

	utility.SetFlash(w, "error", []byte(invalidError))
	http.Redirect(w, r, path, 302)
}

func IsAuthenticated(r *http.Request) bool {
	self := r.Context().Value(middleware.UserContextKey)
	if self != nil && self.(structs.User).ID > 0 {
		return true
	}
	return false
}
