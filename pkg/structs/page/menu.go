package page

// Menu is a component but requires a slice of menu items
type Menu interface {
	GetIdentifier() string
	GetItems() []MenuItem
	SetItems([]MenuItem)
	Render() string
}

type MenuItem struct {
	Text string
	Uri  string
}
