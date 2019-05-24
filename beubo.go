package beubo

type Domain struct {
	Name string
}

var domainList []Domain

func Init() {
	settingsInit()
	databaseInit()
	databaseSeed()
	routesInit()
}
