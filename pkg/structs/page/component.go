package page

type Component interface {
	// render returns a html template string with the content of the field
	Render() string
}
