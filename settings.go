package beubo

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/joho/godotenv"
	"github.com/uberswe/beubo/pkg/routes"
	"github.com/uberswe/beubo/pkg/structs"
	"github.com/uberswe/beubo/pkg/structs/page"
	"github.com/uberswe/beubo/pkg/structs/page/menu"
	"github.com/uberswe/beubo/pkg/template"
	"github.com/uberswe/beubo/pkg/utility"
	"github.com/urfave/negroni"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	databaseHost     = "localhost"
	databaseName     = ""
	databaseUser     = ""
	databasePassword = ""
	databasePort     = "3306"
	databaseDriver   = "mysql"
	environment      = "production"
	testuser         = ""
	testpass         = ""

	rootDir         = "./themes/"
	currentTheme    = "default"
	installed       = false // TODO handle this in a middleware or something
	reloadTemplates = false

	sessionKey = securecookie.GenerateRandomKey(64)

	failures map[string]map[string]string
)

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Printf("[beubo] %s | %s", time.Now().Format("2006-01-02T15:04:05-07:00"), string(bytes))
} //2006-01-02T15:04:05.999Z

// Initializes the settings and checks if Beubo is installed by checking the env file
// If no env file is present then this function will start it's own http listener to go through the installation process
func settingsInit() {

	log.SetFlags(0)
	log.SetOutput(new(logWriter))

	err := godotenv.Load()

	if err != nil {
		// No .env file
		utility.ErrorHandler(err, false)
		log.Println("Attempting to create .env file")
		writeEnv("", "", "", "", "", "", "")
	}

	rootDir = setSetting(os.Getenv("ASSETS_DIR"), rootDir)
	currentTheme = setSetting(os.Getenv("THEME"), currentTheme)

	databaseHost = setSetting(os.Getenv("DB_HOST"), databaseHost)
	databaseName = setSetting(os.Getenv("DB_NAME"), databaseName)
	databaseUser = setSetting(os.Getenv("DB_USER"), databaseUser)
	databaseDriver = setSetting(os.Getenv("DB_DRIVER"), databaseDriver)
	databasePassword = setSetting(os.Getenv("DB_PASSWORD"), databasePassword)
	shouldRefreshDatabase = setBoolSetting(os.Getenv("REFRESH_DATABASE"), shouldRefreshDatabase)
	shouldSeed = setBoolSetting(os.Getenv("SEED_DATABASE"), shouldSeed)

	environment = setSetting(os.Getenv("ENVIRONMENT"), environment)

	testpass = setSetting(os.Getenv("TEST_PASS"), testpass)
	testuser = setSetting(os.Getenv("TEST_USER"), testuser)

	sessionKeyBytes, err := hex.DecodeString(os.Getenv("SESSION_KEY"))
	utility.ErrorHandler(err, false)
	if len(sessionKeyBytes) > 1 {
		sessionKey = sessionKeyBytes
	}

	if databaseName != "" {
		installed = true
	} else {
		log.Println("No installation detected, starting install server")
		srv := startInstallServer()

		// TODO there might be a bug here where we might have multiple instances waiting for installed to be true which causes an infinite loop
		// TODO make this use a channel instead of a loop to wait for install to finish
		for !installed {
			// Pause for 100 ms, this was causing high cpu load without this here
			time.Sleep(time.Millisecond * 100)
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

// startInstallServer starts a http listener and presents the installation page to anyone visiting the port on the host
func startInstallServer() *http.Server {
	r := mux.NewRouter()
	n := negroni.Classic()

	beuboTemplateRenderer := template.BeuboTemplateRenderer{
		ReloadTemplates: true,
		CurrentTheme:    "install",
	}

	beuboTemplateRenderer.Init()

	beuboRouter := &routes.BeuboRouter{
		Renderer: &beuboTemplateRenderer,
	}

	r.NotFoundHandler = http.HandlerFunc(beuboRouter.NotFoundHandler)

	log.Println("Registering themes...")

	r = registerStaticFiles(r)

	log.Println("Registering routes...")

	r.HandleFunc("/", Install)

	n.UseHandler(r)

	srv := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: n}

	log.Println("HTTP Server listening on:", port)
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

// Install handles installation requests and presents the install page
func Install(w http.ResponseWriter, r *http.Request) {

	beuboTemplateRenderer := template.BeuboTemplateRenderer{
		ReloadTemplates: true,
		CurrentTheme:    "install",
	}

	beuboTemplateRenderer.Init()

	pageData := structs.PageData{
		Template: "finished",
		Title:    "Install",
	}

	formKey := "form"
	dbhostKey := "dbhost"
	dbnameKey := "dbname"
	dbuserKey := "dbuser"
	dbpasswordKey := "dbpassword"
	dbdriverKey := "dbdriver"
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
			utility.ErrorHandler(err, false)
		}

		extra[formKey][dbhostKey] = r.PostFormValue(dbhostKey)
		extra[formKey][dbnameKey] = r.PostFormValue(dbnameKey)
		extra[formKey][dbuserKey] = r.PostFormValue(dbuserKey)
		extra[formKey][dbpasswordKey] = r.PostFormValue(dbpasswordKey)
		extra[formKey][dbdriverKey] = r.PostFormValue(dbdriverKey)
		extra[formKey][usernameKey] = r.PostFormValue(usernameKey)
		extra[formKey][passwordKey] = r.PostFormValue(passwordKey)

		token, err := utility.GenerateToken(30)
		if err != nil {
			panic(err)
		}

		failures[token] = extra[formKey]

		utility.SetFlash(w, "token", []byte(token))

		if len(extra[formKey][usernameKey]) < 8 || len(extra[formKey][passwordKey]) < 8 {
			err = errors.New("username and password must be filled with a minimum of 8 characters")
			utility.SetFlash(w, "error", []byte(err.Error()))
			// Redirect back with error
			w.Header().Add("Location", "/")
			w.WriteHeader(302)
			return
		}

		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second, // Slow SQL threshold
				LogLevel:      logger.Info, // Log level
				Colorful:      true,
			},
		)
		config := gorm.Config{
			Logger: newLogger,
		}
		if extra[formKey][dbdriverKey] == "sqlite3" {
			config = gorm.Config{
				DisableForeignKeyConstraintWhenMigrating: true,
				Logger:                                   newLogger,
			}
		}

		dialector := getDialector(extra[formKey][dbuserKey], extra[formKey][dbpasswordKey], extra[formKey][dbhostKey], databasePort, extra[formKey][dbnameKey], extra[formKey][dbdriverKey])
		_, err = gorm.Open(dialector, &config)
		if err != nil {
			utility.SetFlash(w, "error", []byte(err.Error()))
			// Redirect back with error
			w.Header().Add("Location", "/")
			w.WriteHeader(302)
			return
		}

		log.Println("no error, install done")
		writeEnv("", "", extra[formKey][dbhostKey], extra[formKey][dbnameKey], extra[formKey][dbuserKey], extra[formKey][dbpasswordKey], extra[formKey][dbdriverKey])
		beuboTemplateRenderer.RenderHTMLPage(w, r, pageData)
		currentTheme = "default"
		prepareSeed(extra[formKey][usernameKey], extra[formKey][passwordKey])
		databaseName = extra[formKey][dbnameKey]
		installed = true
		return
	}

	menuItems := []page.MenuItem{
		{Text: "Home", URI: "/"},
	}

	menus := []page.Menu{menu.DefaultMenu{
		Items:      menuItems,
		Identifier: "header",
		T:          beuboTemplateRenderer.T,
	}}

	extra = make(map[string]map[string]string)
	token, err := utility.GetFlash(w, r, "token")
	if err == nil {
		extra[formKey] = make(map[string]string)
		extra[formKey] = failures[string(token)]
		failures[string(token)] = nil
	}
	pageData = structs.PageData{
		Theme:    "install",
		Template: "page",
		Title:    "Install",
		Extra:    extra,
		Menus:    menus,
	}
	beuboTemplateRenderer.RenderHTMLPage(w, r, pageData)
	return

}

// setSetting returns the key value if it is set and otherwise falls back to return the variable value
func setSetting(key string, variable string) string {
	if key != "" {
		variable = key
	}
	return variable
}

func setBoolSetting(key string, variable bool) bool {
	if key == "true" {
		return true
	} else if key == "false" {
		return false
	}
	return variable
}

// writeEnv writes environmental variables to an .env file
func writeEnv(assetDir string, theme string, dbHost string, dbName string, dbUser string, dbPassword string, dbDriver string) {
	envContent := []byte("ASSETS_DIR=" + assetDir + "\nTHEME=" + theme + "\n\nDB_DRIVER=" + dbDriver + "\nDB_HOST=" + dbHost + "\nDB_NAME=" + dbName + "\nDB_USER=" + dbUser + "\nDB_PASSWORD=" + dbPassword + "\nSESSION_KEY=" + hex.EncodeToString(sessionKey))
	// TODO allow users to specify folder or even config filename, maybe beuboConfig
	err := ioutil.WriteFile(".env", envContent, 0600) // TODO allow user to change permissions here?

	// We panic if we can not write env
	utility.ErrorHandler(err, false)
}
