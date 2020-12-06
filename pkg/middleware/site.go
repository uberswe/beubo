package middleware

import (
	"github.com/uberswe/beubo/pkg/structs"
	"golang.org/x/net/context"
	"net/http"
)

// Site determines if the domain is an existing site and performs relevant actions based on this
func (bmw *BeuboMiddleware) Site(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	site := structs.FetchSiteByHost(bmw.DB, r.Host)
	if site.ID == 0 {
		// No site detected
		// TODO maybe we should redirect or something if this is the case? Make it configurable
	} else {
		// Site exists
		ctx := context.WithValue(r.Context(), SiteContextKey, site)
		r = r.WithContext(ctx)
	}
	next(rw, r)
}
