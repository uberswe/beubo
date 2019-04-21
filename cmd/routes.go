package cmd

import (
	"encoding/json"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/gorilla/mux"
	"github.com/markustenghamn/beubo/cmd/models"
	"github.com/urfave/negroni"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

var tmpl *template.Template
var rootDir = "web/static/"
var currentTheme = ""
var themes []string
var fileServers = map[string]http.Handler{}

type PageData struct {
	Title       string
	WebsiteName string
	Url         string
	Menu        []MenuItem
	Error       string
	Warning     string
	Message     string
	Year        string
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
	fontFs := http.FileServer(http.Dir("web/static/fonts/"))

	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", cssFs))
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", jsFs))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imgFs))
	r.PathPrefix("/favicon.ico").Handler(imgFs)
	r.PathPrefix("/fonts/").Handler(http.StripPrefix("/fonts/", fontFs))

	log.Println("Registering themes...")

	files, err := ioutil.ReadDir("web/static/themes/")
	checkErr(err)
	for _, f := range files {
		themes = append(themes, f.Name())
		// Register file paths for themes
		fileServers[f.Name()+"_css"] = http.FileServer(http.Dir("web/static/themes/" + f.Name() + "/css/"))
		fileServers[f.Name()+"_js"] = http.FileServer(http.Dir("web/static/themes/" + f.Name() + "/js/"))
		fileServers[f.Name()+"_images"] = http.FileServer(http.Dir("web/static/themes/" + f.Name() + "/images/"))
		fileServers[f.Name()+"_fonts"] = http.FileServer(http.Dir("web/static/themes/" + f.Name() + "/fonts/"))

		r.PathPrefix("/" + f.Name() + "/css/").Handler(http.StripPrefix("/"+f.Name()+"/css/", fileServers[f.Name()+"_css"]))
		r.PathPrefix("/" + f.Name() + "/js/").Handler(http.StripPrefix("/"+f.Name()+"/js/", fileServers[f.Name()+"_js"]))
		r.PathPrefix("/" + f.Name() + "/images/").Handler(http.StripPrefix("/"+f.Name()+"/images/", fileServers[f.Name()+"_images"]))
		r.PathPrefix("/" + f.Name() + "/favicon.ico").Handler(fileServers["/"+f.Name()+"_images"])
		r.PathPrefix("/" + f.Name() + "/fonts/").Handler(http.StripPrefix("/"+f.Name()+"/fonts/", fileServers[f.Name()+"_fonts"]))
	}

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

	// Loads theme templates if defined and falls back to base otherwise
	if currentTheme != "" && tmpl.Lookup(currentTheme+"."+pageTemplate) != nil {
		pageTemplate = currentTheme + "." + pageTemplate
	}

	// Session flash messages to prompt failed logins etc..
	errorMessage, err := GetFlash(w, r, "error")
	errHandler(err)
	warningMessage, err := GetFlash(w, r, "warning")
	errHandler(err)
	stringMessage, err := GetFlash(w, r, "message")
	errHandler(err)

	data := PageData{
		Title:       pageTitle,
		WebsiteName: "qby.se",
		Url:         "http://localhost:3000",
		Menu: []MenuItem{
			{Title: "Home", Path: "/"},
			{Title: "Register", Path: "/register"},
			{Title: "Login", Path: "/login"},
		},
		Error:   string(errorMessage),
		Warning: string(warningMessage),
		Message: string(stringMessage),
		Year:    strconv.Itoa(time.Now().Year()),
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
