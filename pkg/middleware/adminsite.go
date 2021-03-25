package middleware

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/uberswe/beubo/pkg/structs"
	"net/http"
	"strconv"
	"strings"
)

// Site determines if the domain is an existing site and performs relevant actions based on this
func (bmw *BeuboMiddleware) AdminSite(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "admin/sites/a/") {
			params := mux.Vars(r)
			siteID := params["id"]
			siteIDInt, err := strconv.Atoi(siteID)
			if err == nil {
				site := structs.FetchSite(bmw.DB, siteIDInt)
				if site.ID != 0 {
					// Site exists
					ctx := context.WithValue(r.Context(), AdminSiteContextKey, site)
					r = r.WithContext(ctx)
				}
			}

		}
		next.ServeHTTP(w, r)
	})
}
