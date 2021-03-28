package routes

import (
	"github.com/uberswe/beubo/pkg/middleware"
	"github.com/uberswe/beubo/pkg/structs"
	"github.com/uberswe/beubo/pkg/utility"
	"net/http"
	"strconv"
)

func currentUserCanAccessSite(siteId string, br *BeuboRouter, r *http.Request) bool {
	i, err := strconv.Atoi(siteId)

	utility.ErrorHandler(err, false)

	site := structs.FetchSite(br.DB, i)

	self := r.Context().Value(middleware.UserContextKey)
	if self != nil && self.(structs.User).ID > 0 {
		if !self.(structs.User).CanAccessSite(br.DB, site) {
			return false
		}
	}
	return true
}
