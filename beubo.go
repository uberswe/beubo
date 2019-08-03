package beubo

// Init is called to start Beubo, this calls various other functions that initialises
// the database, settings and routes for example.
func Init() {
	settingsInit()
	databaseInit()
	databaseSeed()
	routesInit()
}
