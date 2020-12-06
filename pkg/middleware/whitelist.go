package middleware

import (
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/stretchr/stew/slice"
	"github.com/uberswe/beubo/pkg/structs"
	"log"
	"net"
	"net/http"
	"strings"
)

// Whitelist checks the ip whitelist configuration. If ip whitelisting is enabled, it ensures the ip is whitelisted when accessing administrator pages
func (bmw *BeuboMiddleware) Whitelist(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	settings := structs.FetchSettings(bmw.DB)

	whitelistEnabled := false
	var whitelistedIPs []string

	for _, s := range settings {
		if s.Key == "ip_whitelist" && s.Value == "true" {
			whitelistEnabled = true
		}
		if s.Key == "whitelisted_ip" {
			whitelistedIPs = append(whitelistedIPs, s.Value)
		}
	}

	// TODO if we are behind a load balancer we need to support X-Forwarded-For headers
	// TODO in the future the admin path should be configurable and also not apply to every website
	if strings.HasPrefix(r.URL.Path, "/admin") && whitelistEnabled && !slice.Contains(whitelistedIPs, hostWithoutPort(r.RemoteAddr)) {
		session := sessions.GetSession(r)
		session.Delete("SES_ID")
		session.Clear()
		log.Printf("IP %s is not whitelisted, redirect to /login\n", hostWithoutPort(r.RemoteAddr))
		http.Redirect(rw, r, "/login", 302)
		return
	}
	next(rw, r)
}

func hostWithoutPort(host string) string {
	if strings.Contains(host, ":") {
		h, _, err := net.SplitHostPort(host)
		if err == nil {
			return h
		}
	}
	return host
}
