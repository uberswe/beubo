package plugin

import (
	"net/http"
)

// BeforeRequest is called by Beubo as early as possible in the request stack after other middlewares have executed
func (p Handler) BeforeRequest(w http.ResponseWriter, r *http.Request) {
	for _, plug := range p.Plugins {
		if p.isActive(r, plug) {
			symbol, err := plug.Plugin.Lookup("BeforeRequest")
			if err != nil {
				continue
			}

			beforeRequestFunc, ok := symbol.(func(http.ResponseWriter, *http.Request))
			if !ok {
				continue
			}
			beforeRequestFunc(w, r)
		}
	}
}

// AfterRequest is called by Beubo at the last possible place in the request stack before returning the response writer
func (p Handler) AfterRequest(w http.ResponseWriter, r *http.Request) {
	for _, plug := range p.Plugins {
		if p.isActive(r, plug) {
			symbol, err := plug.Plugin.Lookup("AfterRequest")
			if err != nil {
				continue
			}

			afterRequestFunc, ok := symbol.(func(http.ResponseWriter, *http.Request))
			if !ok {
				continue
			}
			afterRequestFunc(w, r)
		}
	}
}

// PageHandler is called when a non-default route is called in Beubo, returning true will prevent any other handler from executing
func (p Handler) PageHandler(w http.ResponseWriter, r *http.Request) (handled bool) {
	for _, plug := range p.Plugins {
		if p.isActive(r, plug) {
			symbol, err := plug.Plugin.Lookup("PageHandler")
			if err != nil {
				continue
			}

			pageHandlerFunc, ok := symbol.(func(http.ResponseWriter, *http.Request) (handled bool))
			if !ok {
				continue
			}
			return pageHandlerFunc(w, r)
		}
	}
	return false
}
