package routes

import (
	"encoding/json"
	"github.com/markustenghamn/beubo/pkg/utility"
	"net/http"
)

// APIHandler is a prototype route for making base API routes
// TODO implement an external API, preferably with concepts taken from Bridgely, see ticket #15
func (br *BeuboRouter) APIHandler(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal("{'API Test':'Works!'}")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err := w.Write(data)
	utility.ErrorHandler(err, false)
}
