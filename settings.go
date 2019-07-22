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
	"time"
)

var (
	// TODO make this port configurable as an argument
	port = ":3000"

	databaseHost     = "localhost"
	databaseName     = ""
	databaseUser     = ""
	databasePassword = ""
	databasePort     = "3306"

	rootDir      = "./web/"
	currentTheme = "install"
	installed    = false // TODO handle this in a middleware or something

	failures map[string]map[string]string
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

		// TODO there might be a bug here where we might have multiple instances waiting for installed to be true
		for !installed {
			// Pause for 100 ms, this was causing high cpu load without this here
			time.Sleep(time.Second / 10)
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
		log.Println("Server stopped")
	}()

	// returning reference so caller can call Shutdown()
	return srv
}

func Install(w http.ResponseWriter, r *http.Request) {
	formKey := "form"
	dbhostKey := "dbhost"
	dbnameKey := "dbname"
	dbuserKey := "dbuser"
	dbpasswordKey := "dbpassword"
	usernameKey := "username"
	passwordKey := "password"

	if failures == nil {
		failures = make(map[string]map[string]string)
	}

	extra := make(map[string]map[string]string)

	if r.Method == http.MethodPost {
		extra[formKey] = make(map[string]string)
		err := r.ParseForm()
		if err != nil {
			errHandler(err)
		}

		extra[formKey][dbhostKey] = r.PostFormValue(dbhostKey)
		extra[formKey][dbnameKey] = r.PostFormValue(dbnameKey)
		extra[formKey][dbuserKey] = r.PostFormValue(dbuserKey)
		extra[formKey][dbpasswordKey] = r.PostFormValue(dbpasswordKey)
		extra[formKey][usernameKey] = r.PostFormValue(usernameKey)
		extra[formKey][passwordKey] = r.PostFormValue(passwordKey)

		token, err := generateToken(30)
		if err != nil {
			panic(err)
		}

		failures[token] = extra[formKey]

		SetFlash(w, "token", []byte(token))

		if len(extra[formKey][usernameKey]) < 8 || len(extra[formKey][passwordKey]) < 8 {
			err = errors.New("username and password must be filled with a minimum of 8 characters")
			SetFlash(w, "error", []byte(err.Error()))
			// Redirect back with error
			w.Header().Add("Location", "/")
			w.WriteHeader(302)
			return
		}

		connectString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", extra[formKey][dbuserKey], extra[formKey][dbpasswordKey], extra[formKey][dbhostKey], databasePort, extra[formKey][dbnameKey])

		db, err := gorm.Open("mysql", connectString)
		if err != nil {

			SetFlash(w, "error", []byte(err.Error()))
			// Redirect back with error
			w.Header().Add("Location", "/")
			w.WriteHeader(302)
			return
		} else {
			fmt.Println("no error, install done")
			err2 := db.Close()
			errHandler(err2)
			writeEnv("", "", extra[formKey][dbhostKey], extra[formKey][dbnameKey], extra[formKey][dbuserKey], extra[formKey][dbpasswordKey])
			renderHtmlPage("Install", "finished", w, r, nil)
			currentTheme = "default"
			prepareSeed(extra[formKey][usernameKey], extra[formKey][passwordKey])
			// TODO save path to database, could be moved to seed
			paths = append(paths, Path{String: "/admin"})
			installed = true
			return
		}

	} else {
		extra := make(map[string]map[string]string)
		token, err := GetFlash(w, r, "token")
		if err == nil {
			extra[formKey] = make(map[string]string)
			extra[formKey] = failures[string(token)]
			failures[string(token)] = nil
		}
		renderHtmlPage("Install", "page", w, r, extra)
		return
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
