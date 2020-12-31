package plugin

import (
	"github.com/uberswe/beubo/pkg/structs"
	"net/http"
)

func (p Handler) isActive(r *http.Request, plugin Plugin) bool {
	site := structs.FetchSiteByHost(p.DB, r.Host)
	ps := FetchPluginSites(p.DB, plugin.Definition)
	if site.ID >= 0 {
		for _, pluginSite := range ps {
			if pluginSite.SiteID == site.ID && pluginSite.Active {
				return true
			}
		}
	}
	return false
}
