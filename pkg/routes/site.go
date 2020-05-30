package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/markustenghamn/beubo/pkg/structs"
	"github.com/markustenghamn/beubo/pkg/utility"
	"net/http"
	"strconv"
	"strings"
)

func (br *BeuboRouter) SiteAdmin(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	i, err := strconv.Atoi(id)

	utility.ErrorHandler(err, false)

	site := structs.FetchSite(br.DB, i)

	var pages []structs.Page

	extra := make(map[string]interface{})
	pagesRes := make(map[string]map[string]string)
	extra["SiteID"] = fmt.Sprintf("%d", site.ID)

	if err := br.DB.Find(&pages).Error; err != nil {
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

	br.Renderer.RenderHTMLPage("Admin", "admin.site.page.home", w, r, extra)
}

// AdminSiteAdd is the route for adding a site
func (br *BeuboRouter) AdminSiteAdd(w http.ResponseWriter, r *http.Request) {
	br.Renderer.RenderHTMLPage("Admin - Add Site", "admin.site.add", w, r, nil)
}

// Handles adding of a site
func (br *BeuboRouter) AdminSiteAddPost(w http.ResponseWriter, r *http.Request) {
	// TODO should authentication be checked here, maybe with a middleware?

	successMessage := "Site created"
	invalidError := "an error occured and the site could not be created."

	title := r.FormValue("titleField")
	domain := r.FormValue("domainField")
	ssl := r.FormValue("configureSsl")

	domain = strings.ToLower(domain)
	domain = utility.TrimWhitespace(domain)

	if len(title) < 1 {
		invalidError = "The title is too short"
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, "/admin/sites/add", 302)
		return
	}
	if len(domain) < 1 {
		invalidError = "The domain is too short"
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, "/admin/sites/add", 302)
		return
	}

	sslBool := false
	if ssl == "on" {
		sslBool = true
	}

	if structs.CreateSite(br.DB, title, domain, sslBool) {
		utility.SetFlash(w, "message", []byte(successMessage))
		http.Redirect(w, r, "/admin/", 302)
	}

	utility.SetFlash(w, "error", []byte(invalidError))
	http.Redirect(w, r, "/admin/sites/add", 302)
}

func (br *BeuboRouter) AdminSiteDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	i, err := strconv.Atoi(id)

	utility.ErrorHandler(err, false)

	structs.DeleteSite(br.DB, i)

	utility.SetFlash(w, "message", []byte("Site deleted"))

	http.Redirect(w, r, "/admin/", 302)
}

// AdminSiteEdit is the route for adding a site
func (br *BeuboRouter) AdminSiteEdit(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	i, err := strconv.Atoi(id)

	utility.ErrorHandler(err, false)

	site := structs.FetchSite(br.DB, i)

	if site.ID == 0 {
		br.NotFoundHandler(w, r)
		return
	}

	br.Renderer.RenderHTMLPage("Admin - Edit Site", "admin.site.edit", w, r, site)
}

// Handles editing of a site
func (br *BeuboRouter) AdminSiteEditPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	path := fmt.Sprintf("/admin/sites/%s", id)

	i, err := strconv.Atoi(id)

	utility.ErrorHandler(err, false)

	successMessage := "Site updated"
	invalidError := "an error occured and the site could not be updated."

	title := r.FormValue("titleField")
	domain := r.FormValue("domainField")
	ssl := r.FormValue("configureSsl")

	domain = strings.ToLower(domain)
	domain = utility.TrimWhitespace(domain)

	if len(title) < 1 {
		invalidError = "The title is too short"
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, path, 302)
		return
	}
	if len(domain) < 1 {
		invalidError = "The domain is too short"
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, path, 302)
		return
	}

	sslBool := false
	if ssl == "on" {
		sslBool = true
	}

	if structs.UpdateSite(br.DB, i, title, domain, sslBool) {
		utility.SetFlash(w, "message", []byte(successMessage))
		http.Redirect(w, r, path, 302)
	}

	utility.SetFlash(w, "error", []byte(invalidError))
	http.Redirect(w, r, path, 302)
}
