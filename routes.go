package beubo

import (
	"encoding/json"
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

type PageData struct {
	Title       string
	WebsiteName string
	Url         string
	Menu        []MenuItem
	Error       string
	Warning     string
	Message     string
	Year        string
	Extra       interface{}
}

type MenuItem struct {
	Title string
	Path  string
}

func routesInit() {
	// TODO make this port configurable as an argument
	var port = ":3000"
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
	err = http.ListenAndServe(port, n)
	if err != nil {
		log.Println(err)
	}
}

func registerStaticFiles(r *mux.Router) *mux.Router {
	var err error

	log.Println("Parsing and loading templates...")
	tmpl, err = findAndParseTemplates(rootDir, template.FuncMap{})

	files, err := ioutil.ReadDir("web/static/themes/")
	checkErr(err)
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		themes = append(themes, f.Name())
		// Register file paths for themes
		fileServers[f.Name()+"_css"] = gzipped.FileServer(http.Dir("web/static/themes/" + f.Name() + "/css/"))
		fileServers[f.Name()+"_js"] = http.FileServer(http.Dir("web/static/themes/" + f.Name() + "/js/"))
		fileServers[f.Name()+"_images"] = http.FileServer(http.Dir("web/static/themes/" + f.Name() + "/images/"))
		fileServers[f.Name()+"_fonts"] = http.FileServer(http.Dir("web/static/themes/" + f.Name() + "/fonts/"))

		r.PathPrefix("/" + f.Name() + "/css/").Handler(http.StripPrefix("/"+f.Name()+"/css/", fileServers[f.Name()+"_css"]))
		r.PathPrefix("/" + f.Name() + "/js/").Handler(http.StripPrefix("/"+f.Name()+"/js/", fileServers[f.Name()+"_js"]))
		r.PathPrefix("/" + f.Name() + "/images/").Handler(http.StripPrefix("/"+f.Name()+"/images/", fileServers[f.Name()+"_images"]))
		r.PathPrefix("/" + f.Name() + "/favicon.ico").Handler(fileServers["/"+f.Name()+"_images"])
		r.PathPrefix("/" + f.Name() + "/fonts/").Handler(http.StripPrefix("/"+f.Name()+"/fonts/", fileServers[f.Name()+"_fonts"]))
	}
	return r
}

func renderHtmlPage(pageTitle string, pageTemplate string, w http.ResponseWriter, r *http.Request, extra interface{}) {

	var err error

	// Loads theme templates if defined and falls back to base otherwise
	if currentTheme != "" && tmpl.Lookup(currentTheme+"."+pageTemplate) != nil {
		pageTemplate = currentTheme + "." + pageTemplate
	}

	// Session flash messages to prompt failed logins etc..
	errorMessage, err := GetFlash(w, r, "error")
	log.Println(string(errorMessage))
	errHandler(err)
	warningMessage, err := GetFlash(w, r, "warning")
	log.Println(string(warningMessage))
	errHandler(err)
	stringMessage, err := GetFlash(w, r, "message")
	log.Println(string(stringMessage))
	errHandler(err)

	log.Println("setting data")

	data := PageData{
		Title:       pageTitle,
		WebsiteName: "Beubo",
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

	log.Println("execute template")

	err = tmpl.ExecuteTemplate(w, pageTemplate, data)
	log.Println("done")
	errHandler(err)
}

func Home(w http.ResponseWriter, r *http.Request) {
	renderHtmlPage("Home", "page", w, r, nil)
}

func Admin(w http.ResponseWriter, r *http.Request) {
	renderHtmlPage("Admin", "admin.home", w, r, nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if !installed {
		_, err := w.Write([]byte("Beubo is not installed"))
		errHandler(err)
		return
	}
	renderHtmlPage("Login", "login", w, r, nil)
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
	renderHtmlPage("Register", "register", w, r, nil)
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
