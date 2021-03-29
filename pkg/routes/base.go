package routes

import (
	"github.com/uberswe/beubo/pkg/structs"
	beuboPage "github.com/uberswe/beubo/pkg/structs/page"
	"github.com/uberswe/beubo/pkg/structs/page/component"
	"html/template"
	"net/http"
)

// NotFoundHandler overrides the default not found handler
func (br *BeuboRouter) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	pageData := structs.PageData{
		Template: "404",
		Title:    "404 Not Found",
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// PageHandler checks if a page exists for the given slug
func (br *BeuboRouter) PageHandler(w http.ResponseWriter, r *http.Request) {
	if !br.PluginHandler.PageHandler(w, r) {
		site := structs.FetchSiteByHost(br.DB, r.Host)
		if site.ID != 0 {
			page := structs.FetchPageBySiteIDAndSlug(br.DB, int(site.ID), r.URL.Path)
			if page.ID != 0 {
				pageData := structs.PageData{
					Template: "page",
					Title:    page.Title,
					// TODO Components should be defined on the page edit page and defined in the db
					Components: []beuboPage.Component{component.Text{
						Content: template.HTML(page.Content),
						Section: "main",
						T:       br.Renderer.T,
					}},
				}

				br.Renderer.RenderHTMLPage(w, r, pageData)
				return
			}
		}
		br.NotFoundHandler(w, r)
	}
}
