package main

import (
	"log"
	"net/http"
)

// Register is called by Beubo to fetch information about the plugin
func Register() map[string]string {
	return map[string]string{
		"name": "Beubo Example Plugin",
		// identifier should be a unique identifier used to differentiate this plugin from other plugins
		"identifier": "beubo_example_plugin",
	}
}

// BeforeRequest is called by Beubo as early as possible in the request stack after other middlewares have executed
func BeforeRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("BeforeRequest called in example plugin")
}

// AfterRequest is called by Beubo at the last possible place in the request stack before returning the response writer
func AfterRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("AfterRequest called in example plugin")
}

// PageHandler is called when a non-default route is called in Beubo, returning true will prevent any other handler from executing
func PageHandler(w http.ResponseWriter, r *http.Request) (handled bool) {
	log.Println("PageHandler called in example plugin")
	return false
}
