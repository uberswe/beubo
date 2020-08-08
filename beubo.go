package beubo

import (
	"flag"
)

var port = 3000

// Init is called to start Beubo, this calls various other functions that initialises some basic settings
func Init() {
	readCLIFlags()
}

// Run runs the main application
func Run() {
	settingsInit()
	loadPlugins()
	databaseInit()
	databaseSeed()
	routesInit()
}

// readCLIFlags parses command line flags such as port number
func readCLIFlags() {
	flag.IntVar(&port, "port", port, "The port you would like the application to listen on")
	flag.Parse()
}
