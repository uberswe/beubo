package cmd

import (
	"database/sql"
	"github.com/goincremental/negroni-sessions"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
)

type PageData struct {
	Title       string
	WebsiteName string
	Url         string
	Menu        []MenuItem
}

type MenuItem struct {
	Title string
	Path  string
}

func Home(w http.ResponseWriter, req *http.Request) {
	var err error

	tmpl := template.Must(template.ParseFiles("web/static/template/base/page.html"))
	data := PageData{
		Title:       "Home",
		WebsiteName: "QBY.se",
		Url:         "http://localhost:3000",
		Menu: []MenuItem{
			{Title: "Home", Path: "/"},
			{Title: "Register", Path: "/register"},
			{Title: "Login", Path: "/login"},
		},
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}

func Login(w http.ResponseWriter, req *http.Request) {
	var err error

	tmpl := template.Must(template.ParseFiles("web/static/template/auth/login.html"))
	data := PageData{
		Title:       "Login",
		WebsiteName: "QBY.se",
		Url:         "http://localhost:3000",
		Menu: []MenuItem{
			{Title: "Home", Path: "/"},
			{Title: "Register", Path: "/register"},
			{Title: "Login", Path: "/login"},
		},
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}

func LoginPost(w http.ResponseWriter, req *http.Request) {
	session := sessions.GetSession(req)

	username := req.FormValue("inputUsername")
	password := req.FormValue("inputPassword")

	var (
		email                string
		password_in_database string
	)

	err := DB.QueryRow("SELECT user_email, user_password FROM users WHERE user_name = $1", username).Scan(&email, &password_in_database)
	if err == sql.ErrNoRows {
		http.Redirect(w, req, "/authfail", 301)
	} else if err != nil {
		log.Print(err)
		http.Redirect(w, req, "/authfail", 301)
	}

	err = bcrypt.CompareHashAndPassword([]byte(password_in_database), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		http.Redirect(w, req, "/authfail", 301)
	} else if err != nil {
		log.Print(err)
		http.Redirect(w, req, "/authfail", 301)
	}

	session.Set("useremail", email)
	http.Redirect(w, req, "/home", 302)
}
