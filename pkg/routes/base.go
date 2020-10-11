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
	// TODO check if url matches a slug
	site := structs.FetchSiteByHost(br.DB, r.Host)
	if site.ID != 0 {
		page := structs.FetchPageBySiteIDAndSlug(br.DB, int(site.ID), r.URL.Path)
		// TODO add ability to select a page template
		// TODO how does the page know what content to render?
		if page.ID != 0 {
			pageData := structs.PageData{
				Template: "page",
				Title:    page.Title,
				Components: []beuboPage.Component{component.Text{
					Content:  template.HTML(page.Content),
					Theme:    "",
					Template: "",
					Class:    "",
					Renderer: *br.Renderer,
				}},
			}

			br.Renderer.RenderHTMLPage(w, r, pageData)
			return
		}
	}

	br.NotFoundHandler(w, r)
}
