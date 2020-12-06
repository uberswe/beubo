package page

// Menu is a component but requires a slice of menu items
type Menu interface {
	GetIdentifier() string
	GetItems() []MenuItem
	SetItems([]MenuItem)
	Render() string
}

// MenuItem is part of a Menu and usually represents a clickable link
type MenuItem struct {
	Text string
	URI  string
}
