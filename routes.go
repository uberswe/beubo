package beubo

import (
	"encoding/json"
	"fmt"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/gorilla/mux"
	"github.com/lpar/gzipped"
	"github.com/markustenghamn/beubo/pkg/models"
	"github.com/urfave/negroni"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

var tmpl *template.Template
var themes []string
var fileServers = map[string]http.Handler{}

// PageData is a general structure that holds all data that can be displayed on a page
// using go html templates
type PageData struct {
	Title       string
	WebsiteName string
	URL         string
	Menu        []MenuItem
	Error       string
	Warning     string
	Message     string
	Year        string
	Extra       interface{}
}

// MenuItem is one item that can be part of a nav in the frontend
// TODO might be too specific consider removing/redoing this
type MenuItem struct {
	Title string
	Path  string
}

// routesInit initializes the routes and starts a web server that listens on the specified port
func routesInit() {
	// TODO make this port configurable as an argument
	var err error

	errHandler(err)

	r := mux.NewRouter()
	n := negroni.Classic()

	store := cookiestore.New([]byte("kd8ekdleodjfiek"))
	n.Use(sessions.Sessions("global_session_store", store))

	r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	log.Println("Registering themes...")

	r = registerStaticFiles(r)

	log.Println("Registering routes...")

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
	r.HandleFunc("/admin", Admin)
	r.HandleFunc("/logout", Logout)
	r.HandleFunc("/api", APIHandler)

	n.UseHandler(r)

	log.Println("listening on:", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), n)
	if err != nil {
		log.Println(err)
	}
}

// registerStaticFiles handles the loading of all static files for all templates
func registerStaticFiles(r *mux.Router) *mux.Router {
	var err error

	log.Println("Parsing and loading templates...")
	tmpl, err = findAndParseTemplates(rootDir, template.FuncMap{})

	files, err := ioutil.ReadDir("web/themes/")
	checkErr(err)
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		themes = append(themes, f.Name())
		// Register file paths for themes
		fileServers[f.Name()+"_css"] = gzipped.FileServer(http.Dir("web/themes/" + f.Name() + "/css/"))
		fileServers[f.Name()+"_js"] = http.FileServer(http.Dir("web/themes/" + f.Name() + "/js/"))
		fileServers[f.Name()+"_images"] = http.FileServer(http.Dir("web/themes/" + f.Name() + "/images/"))
		fileServers[f.Name()+"_fonts"] = http.FileServer(http.Dir("web/themes/" + f.Name() + "/fonts/"))

		r.PathPrefix("/" + f.Name() + "/css/").Handler(http.StripPrefix("/"+f.Name()+"/css/", fileServers[f.Name()+"_css"]))
		r.PathPrefix("/" + f.Name() + "/js/").Handler(http.StripPrefix("/"+f.Name()+"/js/", fileServers[f.Name()+"_js"]))
		r.PathPrefix("/" + f.Name() + "/images/").Handler(http.StripPrefix("/"+f.Name()+"/images/", fileServers[f.Name()+"_images"]))
		r.PathPrefix("/" + f.Name() + "/favicon.ico").Handler(fileServers["/"+f.Name()+"_images"])
		r.PathPrefix("/" + f.Name() + "/fonts/").Handler(http.StripPrefix("/"+f.Name()+"/fonts/", fileServers[f.Name()+"_fonts"]))
	}
	return r
}

// renderHTMLPage handles rendering of the html template and should be the last function called before returning the response
func renderHTMLPage(pageTitle string, pageTemplate string, w http.ResponseWriter, r *http.Request, extra map[string]map[string]string) {

	log.Println(r.Host)

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
		WebsiteName: "Beubo",
		URL:         "http://localhost:3000",
		Menu: []MenuItem{
			{Title: "Home", Path: "/"},
			{Title: "Register", Path: "/register"},
			{Title: "Login", Path: "/login"},
		},
		Error:   string(errorMessage),
		Warning: string(warningMessage),
		Message: string(stringMessage),
		Year:    strconv.Itoa(time.Now().Year()),
		Extra:   extra,
	}

	err = tmpl.ExecuteTemplate(w, pageTemplate, data)
	errHandler(err)
}

// Home is the default home route and template
func Home(w http.ResponseWriter, r *http.Request) {
	renderHTMLPage("Home", "page", w, r, nil)
}

// Admin is the default admin route and template
func Admin(w http.ResponseWriter, r *http.Request) {
	renderHTMLPage("Admin", "admin.home", w, r, nil)
}

// Login is the default login route
func Login(w http.ResponseWriter, r *http.Request) {
	if !installed {
		_, err := w.Write([]byte("Beubo is not installed"))
		errHandler(err)
		return
	}
	renderHTMLPage("Login", "login", w, r, nil)
}

// LoginPost handles authentication via post request and verifies a username/password via the database
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

// Register renders the default registration page
func Register(w http.ResponseWriter, r *http.Request) {
	renderHTMLPage("Register", "register", w, r, nil)
}

// RegisterPost handles a registration request and inserts the user into the database
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

// Logout handles a GET logout request and removes the user session
func Logout(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	session.Delete("useremail")
	http.Redirect(w, r, "/", 302)
}

// APIHandler is a prototype route for making base API routes
// TODO implement an external API, preferably with concepts taken from Bridgely, see ticket #15
func APIHandler(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal("{'API Test':'Works!'}")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err := w.Write(data)
	errHandler(err)
}
