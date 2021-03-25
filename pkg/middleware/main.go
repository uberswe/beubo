package middleware

import (
	"github.com/uberswe/beubo/pkg/plugin"
	"gorm.io/gorm"
)

type key string

const (
	// UserContextKey is used to fetch and store the user struct from context
	UserContextKey key = "user"
	// SiteContextKey is used to fetch and store the site struct from context
	SiteContextKey key = "site"
	// AdminSiteContextKey is used to fetch and store the site struct from context, this is the id when managing a site in the admin area
	AdminSiteContextKey key = "adminSite"
)

// BeuboMiddleware holds parameters relevant to Beubo middlewares
type BeuboMiddleware struct {
	DB            *gorm.DB
	PluginHandler *plugin.Handler
}
