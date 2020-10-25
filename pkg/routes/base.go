package routes

import (
	"github.com/markustenghamn/beubo/pkg/structs"
	beuboPage "github.com/markustenghamn/beubo/pkg/structs/page"
	"github.com/markustenghamn/beubo/pkg/structs/page/component"
	"html/template"
	"net/http"
)

// NotFoundHandler overrides the default not found handler
func (br *BeuboRouter) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	pageData := structs.PageData{
		Template: "404",
		Title:    "404 Not Found",
	}

	w.WriteHeader(http.StatusNotFound)

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// PageHandler checks if a page exists for the give slug
func (br *BeuboRouter) PageHandler(w http.ResponseWriter, r *http.Request) {
	site := structs.FetchSiteByHost(br.DB, r.Host)
	if site.ID != 0 {
		page := structs.FetchPageBySiteIDAndSlug(br.DB, int(site.ID), r.URL.Path)
		if page.ID != 0 {
			pageData := structs.PageData{
				Template: "page",
				Title:    page.Title,
				// TODO Components should be defined on the page edit page and defined in the db
				Components: []beuboPage.Component{component.Text{
					Content:  template.HTML(page.Content),
					Theme:    "",
					Template: "",
					Class:    "",
					Section:  "main",
					T:        br.Renderer.T,
				}},
			}

			br.Renderer.RenderHTMLPage(w, r, pageData)
			return
		}
	}

	br.NotFoundHandler(w, r)
}
