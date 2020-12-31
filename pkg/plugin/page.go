package plugin

import (
	"github.com/uberswe/beubo/pkg/structs"
	"net/http"
)

// PageData allows the modification of page data before it is passed to the execute function of the template handler
func (p Handler) PageData(r *http.Request, pd structs.PageData) structs.PageData {
	for _, plug := range p.Plugins {
		if p.isActive(r, plug) {
			symbol, err := plug.Plugin.Lookup("PageData")
			if err != nil {
				continue
			}

			pageDataFunc, ok := symbol.(func(*http.Request, structs.PageData) structs.PageData)
			if !ok {
				continue
			}
			return pageDataFunc(r, pd)
		}
	}
	return pd
}
