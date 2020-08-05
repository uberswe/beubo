package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/markustenghamn/beubo/pkg/structs"
	"github.com/markustenghamn/beubo/pkg/utility"
	"net/http"
	"strconv"
)

func (br *BeuboRouter) SiteAdminPageNew(w http.ResponseWriter, r *http.Request) {
	// TODO get the side id as extra
	params := mux.Vars(r)
	siteID := params["id"]
	extra := map[string]string{
		"SiteID": siteID,
	}
	pageData := structs.PageData{
		Template:  "admin.site.page.add",
		Templates: br.Renderer.GetPageTemplates(),
		Title:     "Admin - Add Page",
		Stylesheets: []string{
			"/default/css/normalize.min.css",
			"/default/css/milligram.min.css",
			"/default/css/style.min.css",
			"/default/css/trumbowyg.min.css",
		},
		Scripts: []string{
			"/default/js/jquery-3.3.1.min.js",
			"/default/js/trumbowyg.min.js",
			"/default/js/wysiwyg.min.js",
		},
		Extra: extra,
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

func (br *BeuboRouter) SiteAdminPageNewPost(w http.ResponseWriter, r *http.Request) {
	// TODO should authentication be checked here, maybe with a middleware?
	params := mux.Vars(r)
	siteID := params["id"]

	successMessage := "Site created"
	invalidError := "an error occured and the site could not be created."

	siteIDInt, err := strconv.ParseInt(siteID, 10, 64)
	if err != nil {
		invalidError = "Invalid site id"
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, "/admin/sites", 302)
		return
	}

	title := r.FormValue("titleField")
	slug := r.FormValue("slugField")
	content := r.FormValue("contentField")
	template := r.FormValue("templateField")

	if len(title) < 1 {
		invalidError = "The title is too short"
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, fmt.Sprintf("/admin/sites/admin/%s/page/new", siteID), 302)
		return
	}

	if structs.CreatePage(br.DB, title, slug, template, content, int(siteIDInt)) {
		utility.SetFlash(w, "message", []byte(successMessage))
		http.Redirect(w, r, fmt.Sprintf("/admin/sites/admin/%s", siteID), 302)
		return
	}

	utility.SetFlash(w, "error", []byte(invalidError))
	http.Redirect(w, r, fmt.Sprintf("/admin/sites/admin/%s/page/new", siteID), 302)
}

func (br *BeuboRouter) AdminSitePageEdit(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	siteID := params["id"]
	pageID := params["pageId"]

	pageIDInt, err := strconv.ParseInt(pageID, 10, 64)
	if err != nil {
		invalidError := "Invalid page id"
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, "/admin/sites", 302)
		return
	}

	siteIDInt, err := strconv.ParseInt(siteID, 10, 64)
	if err != nil {
		invalidError := "Invalid site id"
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, "/admin/sites", 302)
		return
	}

	page := structs.FetchPage(br.DB, int(pageIDInt))

	site := structs.FetchSite(br.DB, int(siteIDInt))

	extra := map[string]string{
		"SiteID":     siteID,
		"PageID":     pageID,
		"Slug":       page.Slug,
		"Title":      page.Title,
		"Content":    page.Content,
		"Template":   page.Template,
		"SiteDomain": site.Domain,
	}

	pageData := structs.PageData{
		Template:  "admin.site.page.edit",
		Templates: br.Renderer.GetPageTemplates(),
		Title:     "Admin - Edit Page",
		Extra:     extra,
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

func (br *BeuboRouter) AdminSitePageEditPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	siteID := params["id"]
	pageID := params["pageId"]

	path := fmt.Sprintf("/admin/sites/admin/%s", siteID)

	i, err := strconv.Atoi(siteID)

	pageIDInt, err := strconv.ParseInt(pageID, 10, 64)
	if err != nil {
		invalidError := "Invalid page id"
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, path, 302)
		return
	}

	utility.ErrorHandler(err, false)

	successMessage := "Page updated"
	invalidError := "an error occured and the page could not be updated."

	title := r.FormValue("titleField")
	slug := r.FormValue("slugField")
	content := r.FormValue("contentField")
	template := r.FormValue("templateField")

	if len(title) < 1 {
		invalidError = "The title is too short"
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, fmt.Sprintf("%s/page/edit/%s", path, siteID), 302)
		return
	}

	if structs.UpdatePage(br.DB, i, title, slug, template, content, int(pageIDInt)) {
		utility.SetFlash(w, "message", []byte(successMessage))
		http.Redirect(w, r, path, 302)
		return
	}

	utility.SetFlash(w, "error", []byte(invalidError))
	http.Redirect(w, r, path, 302)
}

func (br *BeuboRouter) AdminSitePageDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	siteID := params["id"]
	pageID := params["pageId"]

	pageIDInt, err := strconv.Atoi(pageID)

	utility.ErrorHandler(err, false)

	structs.DeletePage(br.DB, pageIDInt)

	utility.SetFlash(w, "message", []byte("Site deleted"))

	http.Redirect(w, r, fmt.Sprintf("/admin/sites/admin/%s", siteID), 302)
}
