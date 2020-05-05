package beubo

import (
	"encoding/json"
	"fmt"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/gorilla/mux"
	"github.com/lpar/gzipped"
	beubo "github.com/markustenghamn/beubo/grpc"
	"github.com/markustenghamn/beubo/pkg/structs"
	"github.com/urfave/negroni"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var tmpl *template.Template
var themes []string
var fileServers = map[string]http.Handler{}
var requestChannel = make(chan beubo.Request)
var responseChannel = make(chan beubo.Response)

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
	var err error

	errHandler(err)

	r := mux.NewRouter()
	n := negroni.Classic()

	store := cookiestore.New([]byte("kd8ekdleodjfiek"))
	n.Use(sessions.Sessions("beubo", store))

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

	admin := r.PathPrefix("/admin").Subrouter()
	admin.HandleFunc("/", Admin)
	admin.HandleFunc("/sites/{id}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		if id == "add" {
			if r.Method == "GET" {
				AdminSiteAdd(w, r)
			} else if r.Method == "POST" {
				AdminSiteAddPost(w, r)
			}
		} else {
			if r.Method == "GET" {
				AdminSiteEdit(w, r)
			} else if r.Method == "POST" {
				AdminSiteEditPost(w, r)
			}
		}
	})

	r.HandleFunc("/logout", Logout)
	r.HandleFunc("/api", APIHandler)

	muxer := http.NewServeMux()
	muxer.Handle("/", r)
	muxer.Handle("/admin/", negroni.New(
		negroni.HandlerFunc(auth),
		negroni.Wrap(r),
	))

	n.UseHandler(muxer)

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
	funcMap := buildFuncMap()
	tmpl, err = findAndParseTemplates(rootDir, funcMap)

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
func renderHTMLPage(pageTitle string, pageTemplate string, w http.ResponseWriter, r *http.Request, extra interface{}) {

	var foundTemplate *template.Template
	//log.Println(r.Host)

	var err error

	// Loads theme templates if defined and falls back to base otherwise
	if currentTheme != "" {
		path := fmt.Sprintf("%s.%s", currentTheme, pageTemplate)
		if foundTemplate = tmpl.Lookup(path); foundTemplate == nil {
			log.Printf("Theme file not found %s\n", path)
			return
		}
	} else {
		log.Println("Theme is not set")
		return
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

	if err != nil {
		log.Println("Could not serialize plugin message")
		return
	}

	err = foundTemplate.Execute(w, data)
	errHandler(err)
}

// Home is the default home route and template
func Home(w http.ResponseWriter, r *http.Request) {
	renderHTMLPage("Home", "page", w, r, nil)
}

// Admin is the default admin route and template
func Admin(w http.ResponseWriter, r *http.Request) {
	var sites []structs.Site

	extra := make(map[string]map[string]map[string]string)
	extra["sites"] = make(map[string]map[string]string)

	if err := DB.Find(&sites).Error; err != nil {
		errHandler(err)
	}

	for _, site := range sites {
		sid := fmt.Sprintf("%d", site.ID)
		extra["sites"][sid] = make(map[string]string)
		extra["sites"][sid]["id"] = sid
		extra["sites"][sid]["title"] = site.Title
		extra["sites"][sid]["domain"] = site.Domain
	}

	renderHTMLPage("Admin", "admin.home", w, r, extra)
}

// AdminSiteAdd is the route for adding a site
func AdminSiteAdd(w http.ResponseWriter, r *http.Request) {
	renderHTMLPage("Admin - Add Site", "admin.site.add", w, r, nil)
}

// Handles adding of a site
func AdminSiteAddPost(w http.ResponseWriter, r *http.Request) {
	// TODO should authentication be checked here, maybe with a middleware?

	successMessage := "Site created"
	invalidError := "an error occured and the site could not be created."

	title := r.FormValue("titleField")
	domain := r.FormValue("domainField")
	ssl := r.FormValue("configureSsl")

	domain = strings.ToLower(domain)
	domain = trimWhitespace(domain)

	if len(title) < 1 {
		invalidError = "The title is too short"
		SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, "/admin/sites/add", 302)
		return
	}
	if len(domain) < 1 {
		invalidError = "The domain is too short"
		SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, "/admin/sites/add", 302)
		return
	}

	sslBool := false
	if ssl == "on" {
		sslBool = true
	}

	if structs.CreateSite(DB, title, domain, sslBool) {
		SetFlash(w, "message", []byte(successMessage))
		http.Redirect(w, r, "/admin/sites/add", 302)
	}

	SetFlash(w, "error", []byte(invalidError))
	http.Redirect(w, r, "/admin/sites/add", 302)
}

// AdminSiteEdit is the route for adding a site
func AdminSiteEdit(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	i, err := strconv.Atoi(id)

	errHandler(err)

	site := structs.FetchSite(DB, i)

	if site.ID == 0 {
		NotFoundHandler(w, r)
		return
	}

	renderHTMLPage("Admin - Edit Site", "admin.site.edit", w, r, site)
}

// Handles editing of a site
func AdminSiteEditPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	path := fmt.Sprintf("/admin/sites/%s", id)

	i, err := strconv.Atoi(id)

	errHandler(err)

	successMessage := "Site updated"
	invalidError := "an error occured and the site could not be updated."

	title := r.FormValue("titleField")
	domain := r.FormValue("domainField")
	ssl := r.FormValue("configureSsl")

	domain = strings.ToLower(domain)
	domain = trimWhitespace(domain)

	if len(title) < 1 {
		invalidError = "The title is too short"
		SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, path, 302)
		return
	}
	if len(domain) < 1 {
		invalidError = "The domain is too short"
		SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, path, 302)
		return
	}

	sslBool := false
	if ssl == "on" {
		sslBool = true
	}

	if structs.UpdateSite(DB, i, title, domain, sslBool) {
		SetFlash(w, "message", []byte(successMessage))
		http.Redirect(w, r, path, 302)
	}

	SetFlash(w, "error", []byte(invalidError))
	http.Redirect(w, r, path, 302)
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

	user := structs.AuthUser(DB, email, password)

	if user == nil || user.ID == 0 {
		SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, "/login", 302)
		return
	}

	session.Set("SES_ID", user.ID)
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

	if !structs.CreateUser(DB, email, password) {
		SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, "/login", 302)
	}

	SetFlash(w, "message", []byte("Registration success, please check your email for further instructions"))

	http.Redirect(w, r, "/login", 302)
}

// Logout handles a GET logout request and removes the user session
func Logout(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	session.Delete("SES_ID")
	session.Clear()
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

func auth(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	session := sessions.GetSession(r)
	userId := session.Get("SES_ID")

	uid, err := strconv.Atoi(fmt.Sprintf("%v", userId))

	errHandler(err)

	user := structs.FetchUser(DB, uid)

	if user.ID == 0 {
		log.Println("user is not logged in")
		http.Redirect(rw, r, "/login", 302)
		return
	}

	next(rw, r)
}
