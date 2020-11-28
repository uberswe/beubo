package middleware

import (
	"fmt"
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/uberswe/beubo/pkg/structs"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"strings"
)

// Auth checks if a user is authenticated and performs redirects if needed. The user struct is set to the request context if authenticated.
func (bmw *BeuboMiddleware) Auth(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	session := sessions.GetSession(r)
	token := session.Get("SES_ID")

	user := structs.FetchUserFromSession(bmw.DB, fmt.Sprintf("%v", token))

	// TODO in the future the admin path should be configurable and also not apply to every website
	if user.ID == 0 && strings.HasPrefix(r.URL.Path, "/admin") {
		log.Println("user is not logged in, redirect to /login")
		http.Redirect(rw, r, "/login", 302)
		return
	} else if user.ID > 0 {
		// If path is login then redirect to /admin
		if strings.HasPrefix(r.URL.Path, "/login") {
			http.Redirect(rw, r, "/admin", 302)
			return
		}
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		r = r.WithContext(ctx)
	}

	next(rw, r)
}
