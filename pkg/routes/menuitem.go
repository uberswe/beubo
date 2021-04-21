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
	"log"
	"net/http"
	"strconv"
	"strings"
)

func buildMenuItemLink(action string, menuID int, siteID int) template.URL {
	link := fmt.Sprintf("/%s", strings.ToLower(action))
	if menuID > 0 {
		link = fmt.Sprintf("%s/%d", link, menuID)
	}
	if siteID > 0 {
		return template.URL(fmt.Sprintf("/admin/sites/a/%d/menus%s", siteID, link))
	}
	return template.URL(fmt.Sprintf("/admin/menus%s", link))
}

// AdminMenuItemAdd is the route for adding a menu item to a menu section or menu item
func (br *BeuboRouter) AdminMenuItemAdd(w http.ResponseWriter, r *http.Request) {
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
		// TODO check if a user can access global menus
	}

	form := component.Form{
		Section: "main",
		Fields: []page.Component{
			component.Text{
				Section: "main",
				Content: "This will add a menu section, this section can later be edited to add menu items.",
				T:       br.Renderer.T,
			},
			component.TextField{
				Label: "Section",
				Name:  "section",
				T:     br.Renderer.T,
			},
			component.TextField{
				Label: "Template",
				Name:  "template",
				Value: "menu.default",
				T:     br.Renderer.T,
			},
			component.Button{
				Section: "main",
				Class:   "btn btn-primary",
				Content: "Add",
				T:       br.Renderer.T,
			},
		},
		T:      br.Renderer.T,
		Method: "POST",
		Action: string(buildMenuItemLink("new", 0, i)),
	}

	tmpl := "admin.page"
	if i > 0 {
		tmpl = "admin.site.page"
	}

	pageData := structs.PageData{
		Template: tmpl,
		Title:    "Admin - Add Menu Item",
		Themes:   br.Renderer.GetThemes(),
		Components: []page.Component{
			form,
		},
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// AdminMenuItemAddPost handles adding of a menuitem
func (br *BeuboRouter) AdminMenuItemAddPost(w http.ResponseWriter, r *http.Request) {
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
		// TODO check if a user can access global menus
	}

	path := buildMenuItemLink("new", 0, i)

	successMessage := "Menu created"
	invalidError := "an error occurred and the menu could not be created."

	section := r.FormValue("section")
	tmpl := r.FormValue("tmpl")

	if len(section) < 1 {
		invalidError = "The section is too short"
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, string(path), 302)
		return
	}

	if menu := structs.CreateMenu(br.DB, section, tmpl, i); menu.ID != 0 {
		utility.SetFlash(w, "message", []byte(successMessage))
		http.Redirect(w, r, string(buildMenuLink("", 0, i)), 302)
		return
	}

	utility.SetFlash(w, "error", []byte(invalidError))
	http.Redirect(w, r, string(path), 302)
}

// AdminMenuItemDelete handles the deletion of a menu item
func (br *BeuboRouter) AdminMenuItemDelete(w http.ResponseWriter, r *http.Request) {
	var err error
	if !middleware.CanAccess(br.DB, "manage_menus", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	params := mux.Vars(r)
	menuId := params["menuId"]
	id := params["id"]
	log.Printf("menu %s %s\n", id, menuId)
	mid, err := strconv.Atoi(menuId)
	utility.ErrorHandler(err, false)
	// global menus have site_id 0
	i := 0
	if id != "" {
		i, err = strconv.Atoi(id)
		utility.ErrorHandler(err, false)

		if !currentUserCanAccessSite(id, br, r) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		// TODO check if a user can access global menus
	}

	log.Printf("Deleted %d %d\n", mid, i)

	structs.DeleteMenu(br.DB, mid, i)

	utility.SetFlash(w, "message", []byte("Menu deleted"))

	http.Redirect(w, r, string(buildMenuItemLink("", 0, i)), 302)
}

// AdminMenuItemEdit is the route for editing a menu item
func (br *BeuboRouter) AdminMenuItemEdit(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_menus", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	params := mux.Vars(r)
	menuId := params["menuId"]
	id := params["id"]
	mid, err := strconv.Atoi(menuId)
	utility.ErrorHandler(err, false)
	// global menus have site_id 0
	i := 0
	if id != "" {
		i, err = strconv.Atoi(id)
		utility.ErrorHandler(err, false)
		if !currentUserCanAccessSite(id, br, r) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		// TODO check if a user can access global menus
	}

	menu := structs.FetchMenu(br.DB, mid)

	if menu.ID == 0 {
		br.NotFoundHandler(w, r)
		return
	}

	extra := make(map[string]interface{})
	extra["BackPath"] = buildMenuItemLink("", 0, i)
	menu.Items = structs.FetchMenuItemsBySectionId(br.DB, int(menu.ID))
	for i, _ := range menu.Items {
		err = br.DB.Model(&menu.Items[i]).Association("Permissions").Find(&menu.Items[i].Permissions)
		utility.ErrorHandler(err, false)
		err = br.DB.Model(&menu.Items[i]).Association("Settings").Find(&menu.Items[i].Settings)
		utility.ErrorHandler(err, false)
		menu.Items[i].T = br.Renderer.T
	}
	extra["Menu"] = menu

	pageData := structs.PageData{
		Template: "admin.menu.edit",
		Title:    "Admin - Edit Menu",
		Extra:    extra,
		Themes:   br.Renderer.GetThemes(),
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// AdminMenuItemEditPost handles editing of a menu item
func (br *BeuboRouter) AdminMenuItemEditPost(w http.ResponseWriter, r *http.Request) {
	var err error
	if !middleware.CanAccess(br.DB, "manage_menus", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	params := mux.Vars(r)
	menuId := params["menuId"]
	id := params["id"]
	mid, err := strconv.Atoi(menuId)
	utility.ErrorHandler(err, false)
	// global menus have site_id 0
	i := 0
	if id != "" {
		i, err = strconv.Atoi(id)
		utility.ErrorHandler(err, false)

		if !currentUserCanAccessSite(id, br, r) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		// TODO check if a user can access global menus
	}

	path := buildMenuItemLink("edit", mid, i)

	utility.ErrorHandler(err, false)

	successMessage := "Menu updated"
	invalidError := "an error occurred and the menu could not be updated."

	section := r.FormValue("sectionField")

	// TODO make rules for models
	if len(section) < 1 {
		invalidError = "The key is too short"
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, string(path), 302)
		return
	}

	if structs.UpdateMenu(br.DB, mid, section) {
		utility.SetFlash(w, "message", []byte(successMessage))
		http.Redirect(w, r, string(buildMenuItemLink("", 0, i)), 302)
		return
	}

	utility.SetFlash(w, "error", []byte(invalidError))
	http.Redirect(w, r, string(path), 302)
}
