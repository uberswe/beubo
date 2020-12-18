package plugin

import "net/http"

func (p Handler) BeforeRequest(w http.ResponseWriter, r *http.Request) {
	for _, plug := range p.Plugins {
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

func (p Handler) AfterRequest(w http.ResponseWriter, r *http.Request) {
	for _, plug := range p.Plugins {
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
