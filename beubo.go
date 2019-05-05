package beubo

func Init() {
	settingsInit()
	databaseInit()
	databaseSeed()
	routesInit()
}
