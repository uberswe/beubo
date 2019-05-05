package beubo

import (
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"os"
)

var (
	databaseHost     = ""
	databaseName     = ""
	databaseUser     = ""
	databasePassword = ""

	rootDir      = "./web/static/"
	currentTheme = "install"
	installed    = false // TODO handle this in a middleware or something
)

func settingsInit() {
	err := godotenv.Load()

	if err != nil {
		// No .env file
		errHandler(err)

		writeEnv("", "", "", "", "", "")
	}

	rootDir = setSetting(os.Getenv("ASSETS_DIR"), rootDir)
	currentTheme = setSetting(os.Getenv("THEME"), currentTheme)

	databaseHost = setSetting(os.Getenv("DB_HOST"), databaseHost)
	databaseName = setSetting(os.Getenv("DB_NAME"), databaseName)
	databaseUser = setSetting(os.Getenv("DB_USER"), databaseUser)
	databasePassword = setSetting(os.Getenv("DB_PASSWORD"), databasePassword)

	if databaseHost != "" && databaseName != "" {
		installed = true
	} else {
		log.Println("Running installation, no database configured")
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
