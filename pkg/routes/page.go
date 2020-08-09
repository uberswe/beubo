package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/markustenghamn/beubo/pkg/structs"
	"github.com/markustenghamn/beubo/pkg/utility"
	"log"
	"net/http"
	"strconv"
)

func (br *BeuboRouter) SiteAdminPageNew(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	siteID := params["id"]
	extra := map[string]string{
		"SiteID": siteID,
	}
	pageData := structs.PageData{
		Template:  "admin.site.page.add",
		Templates: br.Renderer.GetPageTemplates(),
		Title:     "Admin - Add Page",
		Extra:     extra,
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

func (br *BeuboRouter) SiteAdminPageNewPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	siteID := params["id"]

	successMessage := "Page created"
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
	tags := r.FormValue("tagField")
	var tagSlice []structs.Tag
	err = json.Unmarshal([]byte(tags), &tagSlice)
	if err != nil {
		log.Println(err)
	}

	for i, tag := range tagSlice {
		tempTag := structs.Tag{}
		br.DB.Where("value = ?", tag.Value).First(&tempTag)
		if tempTag.ID == 0 && br.DB.NewRecord(tag) {
			br.DB.Create(&tag)
			tagSlice[i] = tag
		} else {
			tagSlice[i] = tempTag
		}
	}

	if len(title) < 1 {
		invalidError = "The title is too short"
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, fmt.Sprintf("/admin/sites/a/%s/page/new", siteID), 302)
		return
	}

	if structs.CreatePage(br.DB, title, slug, tagSlice, template, content, int(siteIDInt)) {
		utility.SetFlash(w, "message", []byte(successMessage))
		http.Redirect(w, r, fmt.Sprintf("/admin/sites/a/%s", siteID), 302)
		return
	}

	utility.SetFlash(w, "error", []byte(invalidError))
	http.Redirect(w, r, fmt.Sprintf("/admin/sites/a/%s/page/new", siteID), 302)
	return
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

	// This should not be a nil slice since we are json encoding it even if it is empty
	tags := []structs.JsonTag{}

	for _, tag := range page.Tags {
		tags = append(tags, structs.JsonTag{
			Value: tag.Value,
		})
	}

	var jsonTags []byte
	jsonTags, err = json.Marshal(tags)
	if err != nil {
		log.Println(err)
	}

	extra := map[string]string{
		"SiteID":     siteID,
		"PageID":     pageID,
		"Slug":       page.Slug,
		"Title":      page.Title,
		"Content":    page.Content,
		"Template":   page.Template,
		"SiteDomain": site.Domain,
		"Tags":       string(jsonTags),
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

	path := fmt.Sprintf("/admin/sites/a/%s", siteID)

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
	tags := r.FormValue("tagField")
	var tagSlice []structs.Tag
	err = json.Unmarshal([]byte(tags), &tagSlice)
	if err != nil {
		log.Println(err)
	}

	for i, tag := range tagSlice {
		tempTag := structs.Tag{}
		br.DB.Where("value = ?", tag.Value).First(&tempTag)
		if tempTag.ID == 0 && br.DB.NewRecord(tag) {
			br.DB.Create(&tag)
			tagSlice[i] = tag
		} else {
			tagSlice[i] = tempTag
		}
	}

	if len(title) < 1 {
		invalidError = "The title is too short"
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, fmt.Sprintf("%s/page/edit/%s", path, siteID), 302)
		return
	}

	if structs.UpdatePage(br.DB, i, title, slug, tagSlice, template, content, int(pageIDInt)) {
		utility.SetFlash(w, "message", []byte(successMessage))
		http.Redirect(w, r, path, 302)
		return
	}

	utility.SetFlash(w, "error", []byte(invalidError))
	http.Redirect(w, r, path, 302)
	return
}

func (br *BeuboRouter) AdminSitePageDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	siteID := params["id"]
	pageID := params["pageId"]

	pageIDInt, err := strconv.Atoi(pageID)

	utility.ErrorHandler(err, false)

	structs.DeletePage(br.DB, pageIDInt)

	utility.SetFlash(w, "message", []byte("page deleted"))

	http.Redirect(w, r, fmt.Sprintf("/admin/sites/a/%s", siteID), 302)
}
