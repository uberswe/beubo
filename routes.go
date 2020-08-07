package beubo

import (
	"fmt"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/gorilla/mux"
	"github.com/lpar/gzipped"
	"github.com/markustenghamn/beubo/pkg/middleware"
	"github.com/markustenghamn/beubo/pkg/routes"
	"github.com/markustenghamn/beubo/pkg/template"
	"github.com/markustenghamn/beubo/pkg/utility"
	"github.com/urfave/negroni"
	"io/ioutil"
	"log"
	"net/http"
)

var themes []string
var fileServers = map[string]http.Handler{}

// routesInit initializes the routes and starts a web server that listens on the specified port
func routesInit() {
	var err error

	beuboTemplateRenderer := template.BeuboTemplateRenderer{
		ReloadTemplates: reloadTemplates,
		CurrentTheme:    currentTheme,
		ThemeDir:        rootDir,
	}
	beuboTemplateRenderer.Init()

	beuboRouter := &routes.BeuboRouter{
		DB:       DB,
		Renderer: &beuboTemplateRenderer,
	}

	beuboMiddleware := &middleware.BeuboMiddleware{DB: DB}

	utility.ErrorHandler(err, true)

	r := mux.NewRouter()
	r.StrictSlash(true)
	n := negroni.Classic()

	store := cookiestore.New([]byte(sessionKey))
	n.Use(sessions.Sessions("beubo", store))

	r.NotFoundHandler = http.HandlerFunc(beuboRouter.PageHandler)

	log.Println("Registering themes...")

	r = registerStaticFiles(r)

	log.Println("Registering routes...")

	r.HandleFunc("/login", beuboRouter.Login).Methods("GET")
	r.HandleFunc("/login", beuboRouter.LoginPost).Methods("POST")

	r.HandleFunc("/register", beuboRouter.Register).Methods("GET")
	r.HandleFunc("/register", beuboRouter.RegisterPost).Methods("POST")

	admin := r.PathPrefix("/admin").Subrouter()
	admin.HandleFunc("/", beuboRouter.Admin)

	admin.HandleFunc("/sites/add", beuboRouter.AdminSiteAdd).Methods("GET")
	admin.HandleFunc("/sites/add", beuboRouter.AdminSiteAddPost).Methods("POST")

	admin.HandleFunc("/sites/delete/{id:[0-9]+}", beuboRouter.AdminSiteDelete)

	admin.HandleFunc("/sites/edit/{id:[0-9]+}", beuboRouter.AdminSiteEdit).Methods("GET")
	admin.HandleFunc("/sites/edit/{id:[0-9]+}", beuboRouter.AdminSiteEditPost).Methods("POST")

	siteAdmin := admin.PathPrefix("/sites/a/{id:[0-9]+}").Subrouter()

	siteAdmin.HandleFunc("/", beuboRouter.SiteAdmin)

	siteAdmin.HandleFunc("/page/new", beuboRouter.SiteAdminPageNew).Methods("GET")
	siteAdmin.HandleFunc("/page/new", beuboRouter.SiteAdminPageNewPost).Methods("POST")

	siteAdmin.HandleFunc("/page/edit/{pageId:[0-9]+}", beuboRouter.AdminSitePageEdit).Methods("GET")
	siteAdmin.HandleFunc("/page/edit/{pageId:[0-9]+}", beuboRouter.AdminSitePageEditPost).Methods("POST")

	siteAdmin.HandleFunc("/page/delete/{pageId:[0-9]+}", beuboRouter.AdminSitePageDelete)

	r.HandleFunc("/logout", beuboRouter.Logout)
	r.HandleFunc("/api", beuboRouter.APIHandler)

	muxer := http.NewServeMux()
	muxer.Handle("/", negroni.New(
		negroni.HandlerFunc(beuboMiddleware.Site),
		negroni.HandlerFunc(beuboMiddleware.Auth),
		negroni.Wrap(r),
	))

	n.UseHandler(muxer)

	log.Println("HTTP Server listening on:", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), n)
	if err != nil {
		log.Println(err)
	}
}

// registerStaticFiles handles the loading of all static files for all templates
func registerStaticFiles(r *mux.Router) *mux.Router {
	var err error

	themedir := "themes/"
	files, err := ioutil.ReadDir(themedir)
	utility.ErrorHandler(err, false)
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		themes = append(themes, f.Name())
		// Register file paths for themes
		fileServers[f.Name()+"_css"] = gzipped.FileServer(http.Dir(themedir + f.Name() + "/css/"))
		fileServers[f.Name()+"_js"] = http.FileServer(http.Dir(themedir + f.Name() + "/js/"))
		fileServers[f.Name()+"_images"] = http.FileServer(http.Dir(themedir + f.Name() + "/images/"))
		fileServers[f.Name()+"_fonts"] = http.FileServer(http.Dir(themedir + f.Name() + "/fonts/"))

		r.PathPrefix("/" + f.Name() + "/css/").Handler(http.StripPrefix("/"+f.Name()+"/css/", fileServers[f.Name()+"_css"]))
		r.PathPrefix("/" + f.Name() + "/js/").Handler(http.StripPrefix("/"+f.Name()+"/js/", fileServers[f.Name()+"_js"]))
		r.PathPrefix("/" + f.Name() + "/images/").Handler(http.StripPrefix("/"+f.Name()+"/images/", fileServers[f.Name()+"_images"]))
		r.PathPrefix("/" + f.Name() + "/favicon.ico").Handler(fileServers["/"+f.Name()+"_images"])
		r.PathPrefix("/" + f.Name() + "/fonts/").Handler(http.StripPrefix("/"+f.Name()+"/fonts/", fileServers[f.Name()+"_fonts"]))
	}
	return r
}
