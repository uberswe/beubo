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

func BeforeRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("BeforeRequest called in example plugin")
}

func AfterRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("AfterRequest called in example plugin")
}
