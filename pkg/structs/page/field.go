package page

type Field interface {
	// render returns a html template string with the content of the field
	Render() string
}
