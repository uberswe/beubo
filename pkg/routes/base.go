package routes

import (
	"net/http"
)

// Home is the default home route and template
func (br *BeuboRouter) Home(w http.ResponseWriter, r *http.Request) {
	br.Renderer.RenderHTMLPage("Home", "page", w, r, nil)
}

// NotFoundHandler overrides the default not found handler
func (br *BeuboRouter) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	br.Renderer.RenderHTMLPage("404 Not Found", "404", w, r, nil)
}
