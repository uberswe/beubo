package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/uberswe/beubo/pkg/middleware"
	"github.com/uberswe/beubo/pkg/structs"
	"github.com/uberswe/beubo/pkg/utility"
	"net/http"
	"strconv"
	"strings"
)

// SiteAdmin is the main page for the admin area and shows a list of pages
func (br *BeuboRouter) SiteAdmin(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_pages", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

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

// AdminSiteAdd is the route for adding a site
func (br *BeuboRouter) AdminSiteAdd(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_sites", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

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

	pageData := structs.PageData{
		Template: "admin.site.add",
		Title:    "Admin - Add Site",
		Themes:   br.Renderer.GetThemes(),
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// AdminSiteAddPost handles the post request for adding a site
func (br *BeuboRouter) AdminSiteAddPost(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_sites", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

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

	path := "/admin/sites/add"

	successMessage := "Site created"
	invalidError := "an error occurred and the site could not be created."

	themeID := 0
	title := r.FormValue("titleField")
	domain := r.FormValue("domainField")
	// typeField
	// 1 - Beubo hosted site
	// 2 - HTML files from directory
	// 3 - redirect to a different domain
	siteType := r.FormValue("typeField")

	typeID, err := strconv.Atoi(siteType)

	utility.ErrorHandler(err, false)
	// Theme is only relevant for Beubo hosted sites
	if siteType == "1" {
		theme := r.FormValue("themeField")
		themeStruct := structs.FetchThemeBySlug(br.DB, theme)
		if themeStruct.ID == 0 {
			invalidError = "The theme is invalid"
			utility.SetFlash(w, "error", []byte(invalidError))
			http.Redirect(w, r, path, 302)
			return
		}
		themeID = int(themeStruct.ID)
	}

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

	if structs.CreateSite(br.DB, title, domain, typeID, themeID) {
		utility.SetFlash(w, "message", []byte(successMessage))
		http.Redirect(w, r, "/admin/", 302)
		return
	}

	utility.SetFlash(w, "error", []byte(invalidError))
	http.Redirect(w, r, "/admin/sites/add", 302)
}

// AdminSiteDelete is the route for deleting a site
func (br *BeuboRouter) AdminSiteDelete(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_sites", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

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

	structs.DeleteSite(br.DB, i)

	utility.SetFlash(w, "message", []byte("Site deleted"))

	http.Redirect(w, r, "/admin/", 302)
}

// AdminSiteEdit is the route for adding a site
func (br *BeuboRouter) AdminSiteEdit(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_sites", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	params := mux.Vars(r)
	id := params["id"]

	i, err := strconv.Atoi(id)

	utility.ErrorHandler(err, false)

	site := structs.FetchSite(br.DB, i)

	if site.ID == 0 {
		br.NotFoundHandler(w, r)
		return
	}

	self := r.Context().Value(middleware.UserContextKey)
	if self != nil && self.(structs.User).ID > 0 {
		if !self.(structs.User).CanAccessSite(br.DB, site) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}

	pageData := structs.PageData{
		Template: "admin.site.edit",
		Title:    "Admin - Edit Site",
		Extra:    site,
		Themes:   br.Renderer.GetThemes(),
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// AdminSiteEditPost handles editing of a site
func (br *BeuboRouter) AdminSiteEditPost(w http.ResponseWriter, r *http.Request) {
	if !middleware.CanAccess(br.DB, "manage_sites", r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	params := mux.Vars(r)
	id := params["id"]

	path := fmt.Sprintf("/admin/sites/edit/%s", id)

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

	successMessage := "Site updated"
	invalidError := "an error occurred and the site could not be updated."

	themeID := 0
	title := r.FormValue("titleField")
	domain := r.FormValue("domainField")
	// typeField
	// 1 - Beubo hosted site
	// 2 - HTML files from directory
	// 3 - redirect to a different domain
	siteType := r.FormValue("typeField")

	typeID, err := strconv.Atoi(siteType)

	utility.ErrorHandler(err, false)
	// Theme is only relevant for Beubo hosted sites
	if siteType == "1" {
		theme := r.FormValue("themeField")
		themeStruct := structs.FetchThemeBySlug(br.DB, theme)
		if themeStruct.ID == 0 {
			invalidError = "The theme is invalid"
			utility.SetFlash(w, "error", []byte(invalidError))
			http.Redirect(w, r, path, 302)
			return
		}
		themeID = int(themeStruct.ID)
	}

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

	if structs.UpdateSite(br.DB, i, title, domain, typeID, themeID) {
		utility.SetFlash(w, "message", []byte(successMessage))
		http.Redirect(w, r, path, 302)
		return
	}

	utility.SetFlash(w, "error", []byte(invalidError))
	http.Redirect(w, r, path, 302)
}
