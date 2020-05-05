package beubo

import "net/http"

// NotFoundHandler overrides the default not found handler
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	renderHTMLPage("404 Not Found", "404", w, r, nil)
}
