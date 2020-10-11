package component

import "github.com/markustenghamn/beubo/pkg/structs/page"

type Table struct {
	// Row defines
	Row Row
}

type Row struct {
	Columns []Column
}

type Column struct {
	Name  string
	Field page.Field
}
