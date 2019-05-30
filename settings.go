package beubo

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/urfave/negroni"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	// TODO make this port configurable as an argument
	port = ":3000"

	databaseHost     = "localhost"
	databaseName     = ""
	databaseUser     = ""
	databasePassword = ""
	databasePort     = "3306"

	rootDir      = "./web/static/"
	currentTheme = "install"
	installed    = false // TODO handle this in a middleware or something
)

func settingsInit() {

	err := godotenv.Load()

	if err != nil {
		// No .env file
		errHandler(err)
		log.Println("Attempting to create .env file")
		writeEnv("", "", "", "", "", "")
	}

	rootDir = setSetting(os.Getenv("ASSETS_DIR"), rootDir)
	currentTheme = setSetting(os.Getenv("THEME"), currentTheme)

	databaseHost = setSetting(os.Getenv("DB_HOST"), databaseHost)
	databaseName = setSetting(os.Getenv("DB_NAME"), databaseName)
	databaseUser = setSetting(os.Getenv("DB_USER"), databaseUser)
	databasePassword = setSetting(os.Getenv("DB_PASSWORD"), databasePassword)

	if databaseUser != "" && databaseName != "" {
		installed = true
		currentTheme = "default"
	} else {
		log.Println("No installation detected, starting install server")
		srv := startInstallServer()

		for !installed {
			// Keep running install server until installed is finished
		}

		if err := srv.Shutdown(context.TODO()); err != nil {
			panic(err) // failure/timeout shutting down the server gracefully
		}
		log.Println("Install complete, restarting server")
		// settingsInit() calls itself after install to reload settings
		settingsInit()
	}

}

func startInstallServer() *http.Server {
	r := mux.NewRouter()
	n := negroni.Classic()

	r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	log.Println("Registering themes...")

	r = registerStaticFiles(r)

	log.Println("Registering routes...")

	r.HandleFunc("/", Install)

	n.UseHandler(r)

	srv := &http.Server{Addr: port, Handler: n}

	log.Println("listening on:", port)
	go func() {
		// returns ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// NOTE: there is a chance that next line won't have time to run,
			// as main() doesn't wait for this goroutine to stop. don't use
			// code with race conditions like these for production. see post
			// comments below on more discussion on how to handle this.
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()

	// returning reference so caller can call Shutdown()
	return srv
}

func Install(w http.ResponseWriter, r *http.Request) {
	// TODO check if installed here, if db is configured

	log.Println(r.Host)

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			errHandler(err)
		}

		domain := r.PostFormValue("domain")
		adminpath := r.PostFormValue("adminpath")
		dbhost := r.PostFormValue("dbhost")
		dbname := r.PostFormValue("dbname")
		dbuser := r.PostFormValue("dbuser")
		dbpassword := r.PostFormValue("dbpassword")
		email := r.PostFormValue("email")
		password := r.PostFormValue("password")

		if len(email) == 0 && len(password) == 0 {
			err = errors.New("email and password must be filled")
		}

		if len(domain) == 0 && len(adminpath) == 0 {
			err = errors.New("email and password must be filled")
		}

		connectString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbuser, dbpassword, dbhost, databasePort, dbname)

		_, err = gorm.Open("mysql", connectString)
		if err != nil {
			// TODO return error and go back to install page
			renderHtmlPage("Install", "page", w, r, nil)
		} else {
			writeEnv("", "", dbhost, dbname, dbuser, dbpassword)
			renderHtmlPage("Install", "finished", w, r, nil)
			currentTheme = "default"
			prepareSeed(email, password)
			// TODO should save these objects to database at some point
			domains = append(domains, Domain{Name: domain})
			paths = append(paths, Path{String: adminpath})
			installed = true
		}

	} else {
		renderHtmlPage("Install", "page", w, r, nil)
	}
}

func setSetting(key string, variable string) string {
	if key != "" {
		variable = key
	}
	return variable
}

func writeEnv(assetDir string, theme string, dbHost string, dbName string, dbUser string, dbPassword string) {
	envContent := []byte("ASSETS_DIR=" + assetDir + "\nTHEME=" + theme + "\n\nDB_HOST=" + dbHost + "\nDB_NAME=" + dbName + "\nDB_USER=" + dbUser + "\nDB_PASSWORD=" + dbPassword)
	// TODO allow users to specify folder or even config filename, maybe beuboConfig
	err := ioutil.WriteFile(".env", envContent, 0600) // TODO allow user to change permissions here?

	// We panic if we can not write env
	checkErr(err)
}
