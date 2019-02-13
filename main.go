package main

import (
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/gorilla/mux"
	"github.com/markustenghamn/beubo/cmd"
	"github.com/urfave/negroni"
	"log"
	"net/http"
)

func main() {
	var port = ":3000"
	var err error

	cmd.Init()

	r := mux.NewRouter()
	n := negroni.Classic()

	store := cookiestore.New([]byte("kd8ekdleodjfiek"))
	n.Use(sessions.Sessions("global_session_store", store))

	r.NotFoundHandler = http.HandlerFunc(cmd.NotFoundHandler)

	cssFs := http.FileServer(http.Dir("web/static/css/"))
	jsFs := http.FileServer(http.Dir("web/static/js/"))
	imgFs := http.FileServer(http.Dir("web/static/images/"))

	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", cssFs))
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", jsFs))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imgFs))
	r.PathPrefix("/favicon.ico").Handler(imgFs)

	r.HandleFunc("/", cmd.Home)

	n.UseHandler(r)

	log.Println("listening on:", port)
	err = http.ListenAndServe(port, n)
	if err != nil {
		log.Println(err)
	}
}
