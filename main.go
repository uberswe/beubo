package main

import (
	"fmt"
	"github.com/urfave/negroni"
	"log"
	"net/http"
)

func main() {
	var err error
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		_, err = fmt.Fprintf(w, "Welcome to the home page!")
		if err != nil {
			log.Println(err)
		}
	})

	n := negroni.Classic() // Includes some default middlewares
	n.UseHandler(mux)

	err = http.ListenAndServe(":3000", n)
	if err != nil {
		log.Println(err)
	}
}
