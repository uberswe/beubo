package middleware

import "net/http"

// Plugin allows plugins to perform actions as a middleware
func (bmw *BeuboMiddleware) Plugin(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	bmw.PluginHandler.BeforeRequest(w, r)
	next.ServeHTTP(w, r)
	bmw.PluginHandler.AfterRequest(w, r)
}
