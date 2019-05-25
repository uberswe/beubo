package beubo

type Domain struct {
	Name string
}

type Path struct {
	String string
}

var (
	domains []Domain
	paths   []Path
)

func Init() {
	settingsInit()
	databaseInit()
	databaseSeed()
	routesInit()
}
