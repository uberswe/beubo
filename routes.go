package beubo

import (
	"fmt"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/gorilla/mux"
	"github.com/uberswe/beubo/pkg/middleware"
	"github.com/uberswe/beubo/pkg/routes"
	"github.com/uberswe/beubo/pkg/template"
	"github.com/uberswe/beubo/pkg/utility"
	"github.com/urfave/negroni"
	"golang.org/x/time/rate"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
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
		PluginHandler:   &pluginHandler,
		DB:              DB,
	}
	beuboTemplateRenderer.Init()

	beuboRouter := &routes.BeuboRouter{
		DB:            DB,
		Renderer:      &beuboTemplateRenderer,
		PluginHandler: &pluginHandler,
	}

	beuboMiddleware := &middleware.BeuboMiddleware{
		DB:            DB,
		PluginHandler: &pluginHandler,
	}

	// TODO make the burst and time configurable
	throttleMiddleware := middleware.Throttle{
		IPs:     make(map[string]middleware.ThrottleClient),
		Rate:    rate.Every(time.Minute),
		Burst:   300,
		Mu:      &sync.RWMutex{},
		Cleanup: time.Duration(24) * time.Hour,
	}

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

	admin.HandleFunc("/settings", beuboRouter.Settings)
	admin.HandleFunc("/users", beuboRouter.Users)
	admin.HandleFunc("/users/roles", beuboRouter.AdminUserRoles)
	admin.HandleFunc("/plugins", beuboRouter.Plugins)

	admin.HandleFunc("/sites/add", beuboRouter.AdminSiteAdd).Methods("GET")
	admin.HandleFunc("/sites/add", beuboRouter.AdminSiteAddPost).Methods("POST")

	admin.HandleFunc("/settings/add", beuboRouter.AdminSettingAdd).Methods("GET")
	admin.HandleFunc("/settings/add", beuboRouter.AdminSettingAddPost).Methods("POST")

	admin.HandleFunc("/users/add", beuboRouter.AdminUserAdd).Methods("GET")
	admin.HandleFunc("/users/add", beuboRouter.AdminUserAddPost).Methods("POST")

	admin.HandleFunc("/users/roles/add", beuboRouter.AdminUserRoleAdd).Methods("GET")
	admin.HandleFunc("/users/roles/add", beuboRouter.AdminUserRoleAddPost).Methods("POST")

	admin.HandleFunc("/sites/delete/{id:[0-9]+}", beuboRouter.AdminSiteDelete)
	admin.HandleFunc("/settings/delete/{id:[0-9]+}", beuboRouter.AdminSettingDelete)
	admin.HandleFunc("/users/delete/{id:[0-9]+}", beuboRouter.AdminUserDelete)
	admin.HandleFunc("/users/roles/delete/{id:[0-9]+}", beuboRouter.AdminUserRoleDelete)

	admin.HandleFunc("/sites/edit/{id:[0-9]+}", beuboRouter.AdminSiteEdit).Methods("GET")
	admin.HandleFunc("/sites/edit/{id:[0-9]+}", beuboRouter.AdminSiteEditPost).Methods("POST")

	admin.HandleFunc("/settings/edit/{id:[0-9]+}", beuboRouter.AdminSettingEdit).Methods("GET")
	admin.HandleFunc("/settings/edit/{id:[0-9]+}", beuboRouter.AdminSettingEditPost).Methods("POST")

	admin.HandleFunc("/users/edit/{id:[0-9]+}", beuboRouter.AdminUserEdit).Methods("GET")
	admin.HandleFunc("/users/edit/{id:[0-9]+}", beuboRouter.AdminUserEditPost).Methods("POST")

	admin.HandleFunc("/users/roles/edit/{id:[0-9]+}", beuboRouter.AdminUserRoleEdit).Methods("GET")
	admin.HandleFunc("/users/roles/edit/{id:[0-9]+}", beuboRouter.AdminUserRoleEditPost).Methods("POST")

	admin.HandleFunc("/plugins/edit/{id:[a-zA-Z_]+}", beuboRouter.AdminPluginEdit).Methods("GET")
	admin.HandleFunc("/plugins/edit/{id:[a-zA-Z_]+}", beuboRouter.AdminPluginEditPost).Methods("POST")

	// TODO I don't like this /sites/a/ structure of the routes, consider changing it
	siteAdmin := admin.PathPrefix("/sites/a/{id:[0-9]+}").Subrouter()

	siteAdmin.HandleFunc("/", beuboRouter.SiteAdmin)

	siteAdmin.HandleFunc("/page/new", beuboRouter.SiteAdminPageNew).Methods("GET")
	siteAdmin.HandleFunc("/page/new", beuboRouter.SiteAdminPageNewPost).Methods("POST")

	siteAdmin.HandleFunc("/page/edit/{pageId:[0-9]+}", beuboRouter.AdminSitePageEdit).Methods("GET")
	siteAdmin.HandleFunc("/page/edit/{pageId:[0-9]+}", beuboRouter.AdminSitePageEditPost).Methods("POST")

	siteAdmin.HandleFunc("/page/delete/{pageId:[0-9]+}", beuboRouter.AdminSitePageDelete)

	r.HandleFunc("/logout", beuboRouter.Logout)

	muxer := http.NewServeMux()
	muxer.Handle("/", negroni.New(
		negroni.HandlerFunc(beuboMiddleware.Site),
		negroni.HandlerFunc(beuboMiddleware.Whitelist),
		negroni.HandlerFunc(beuboMiddleware.Auth),
		negroni.HandlerFunc(throttleMiddleware.Throttle),
		negroni.HandlerFunc(beuboMiddleware.Plugin),
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
		fileServers[f.Name()+"_css"] = http.FileServer(http.Dir(themedir + f.Name() + "/css/"))
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
