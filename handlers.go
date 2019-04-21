package beubo

import "net/http"

// NotFoundHandler overrides the default not found handler
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	// You can use the serve file helper to respond to 404 with
	// your request file.

	http.ServeFile(w, r, "web/static/html/base/404.html")
}
