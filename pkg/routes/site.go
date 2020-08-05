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

	pageData := structs.PageData{
		Template: "admin.site.page.home",
		Title:    "Admin",
		Extra:    extra,
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// AdminSiteAdd is the route for adding a site
func (br *BeuboRouter) AdminSiteAdd(w http.ResponseWriter, r *http.Request) {
	pageData := structs.PageData{
		Template: "admin.site.add",
		Title:    "Admin - Add Site",
		Themes:   br.Renderer.GetThemes(),
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// Handles adding of a site
func (br *BeuboRouter) AdminSiteAddPost(w http.ResponseWriter, r *http.Request) {
	// TODO should authentication be checked here, maybe with a middleware?
	path := "/admin/sites/add"

	successMessage := "Site created"
	invalidError := "an error occured and the site could not be created."

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

	pageData := structs.PageData{
		Template: "admin.site.edit",
		Title:    "Admin - Edit Site",
		Extra:    site,
		Themes:   br.Renderer.GetThemes(),
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// Handles editing of a site
func (br *BeuboRouter) AdminSiteEditPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	path := fmt.Sprintf("/admin/sites/edit/%s", id)

	i, err := strconv.Atoi(id)

	utility.ErrorHandler(err, false)

	successMessage := "Site updated"
	invalidError := "an error occured and the site could not be updated."

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
