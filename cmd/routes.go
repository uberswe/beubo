package cmd

import (
	"encoding/json"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/gorilla/mux"
	"github.com/markustenghamn/beubo/cmd/models"
	"github.com/urfave/negroni"
	"html/template"
	"log"
	"net/http"
)

var tmpl *template.Template
var rootDir = "web/static/template/"

type PageData struct {
	Title       string
	WebsiteName string
	Url         string
	Menu        []MenuItem
	Error       string
	Warning     string
	Message     string
}

type MenuItem struct {
	Title string
	Path  string
}

func InitRoutes() {
	var port = ":3000"
	var err error

	log.Println("Parsing and loading templates...")
	tmpl, err = findAndParseTemplates(rootDir, template.FuncMap{})
	errHandler(err)

	r := mux.NewRouter()
	n := negroni.Classic()

	store := cookiestore.New([]byte("kd8ekdleodjfiek"))
	n.Use(sessions.Sessions("global_session_store", store))

	r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	log.Println("Registering routes...")

	cssFs := http.FileServer(http.Dir("web/static/css/"))
	jsFs := http.FileServer(http.Dir("web/static/js/"))
	imgFs := http.FileServer(http.Dir("web/static/images/"))

	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", cssFs))
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", jsFs))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imgFs))
	r.PathPrefix("/favicon.ico").Handler(imgFs)

	r.HandleFunc("/", Home)
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			Login(w, r)
		} else if r.Method == "POST" {
			LoginPost(w, r)
		}
	})
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			Register(w, r)
		} else if r.Method == "POST" {
			RegisterPost(w, r)
		}
	})
	r.HandleFunc("/logout", Logout)
	r.HandleFunc("/api", APIHandler)

	n.UseHandler(r)

	log.Println("listening on:", port)
	err = http.ListenAndServe(port, n)
	if err != nil {
		log.Println(err)
	}
}

func renderHtmlPage(pageTitle string, pageTemplate string, w http.ResponseWriter, r *http.Request) {

	var err error

	// Session flash messages to prompt failed logins etc..
	errorMessage, err := GetFlash(w, r, "error")
	errHandler(err)
	warningMessage, err := GetFlash(w, r, "warning")
	errHandler(err)
	stringMessage, err := GetFlash(w, r, "message")
	errHandler(err)

	data := PageData{
		Title:       pageTitle,
		WebsiteName: "QBY.se",
		Url:         "http://localhost:3000",
		Menu: []MenuItem{
			{Title: "Home", Path: "/"},
			{Title: "Register", Path: "/register"},
			{Title: "Login", Path: "/login"},
		},
		Error:   string(errorMessage),
		Warning: string(warningMessage),
		Message: string(stringMessage),
	}

	err = tmpl.ExecuteTemplate(w, pageTemplate, data)
	errHandler(err)
}

func Home(w http.ResponseWriter, r *http.Request) {
	renderHtmlPage("Home", "page", w, r)
}

func Login(w http.ResponseWriter, r *http.Request) {
	renderHtmlPage("Login", "login", w, r)
}

func LoginPost(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)

	invalidError := "Email or password incorrect, please try again or contact support"

	email := r.FormValue("email")
	password := r.FormValue("password")
	sessionid := "sessionidgoeshere"

	if !models.AuthUser(DB, email, password) {
		SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, "/login", 302)
		return
	}

	session.Set("SES_ID", sessionid)
	http.Redirect(w, r, "/admin", 302)
}

func Register(w http.ResponseWriter, r *http.Request) {
	renderHtmlPage("Register", "register", w, r)
}

func RegisterPost(w http.ResponseWriter, r *http.Request) {

	invalidError := "Please make sure the email is correct or that it does not already belong to a registered account"

	email := r.FormValue("email")
	password := r.FormValue("password")

	if !models.CreateUser(DB, email, password) {
		SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, "/login", 302)
	}

	SetFlash(w, "message", []byte("Registration success, please check your email for further instructions"))

	http.Redirect(w, r, "/login", 302)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	session.Delete("useremail")
	http.Redirect(w, r, "/", 302)
}

func APIHandler(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal("{'API Test':'Works!'}")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err := w.Write(data)
	errHandler(err)
}
