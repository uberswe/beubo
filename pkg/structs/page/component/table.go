package component

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/markustenghamn/beubo/pkg/structs/page"
	"github.com/markustenghamn/beubo/pkg/utility"
	"html/template"
	"reflect"
)

type Table struct {
	Section          string
	Header           []Column
	Rows             []Row
	Theme            string
	Template         string
	PageNumber       int // Current page
	PageDisplayCount int // How many rows per page
	T                *template.Template
}

func (t Table) GetSection() string {
	return t.Section
}

type Row struct {
	Columns []Column
}

type Column struct {
	Name  string
	Field page.Component
	Value string
}

func (t Table) GetTemplateName() string {
	return returnTIfNotEmpty(t.Template, "component.table")
}

func (t Table) GetTheme() string {
	return returnTIfNotEmpty(t.Template, "default")
}

func (t Table) GetTemplate() *template.Template {
	return t.T
}

func (t Table) Render() string {
	return page.RenderCompnent(t)
}

func (t Table) RenderColumn(c Column) string {
	return page.RenderCompnent(c.Field)
}

func (c Column) RenderField(value string, field page.Component) template.HTML {
	if field != nil && field.Render() != "" {
		return template.HTML(field.Render())
	}
	return template.HTML(value)
}

type ColumnDefinition struct {
	Name                 string
	ValueFromStructField string
	StaticValue          string
	ComponentDefinition  *page.ComponentDefinition
}

func MakeTable(db *gorm.DB, m interface{}, cd []ColumnDefinition, numRows int, page int, section string, theme string, template string, t *template.Template) (Table, error) {
	table := Table{}
	if err := db.Limit(numRows).Offset(page).Find(&m).Error; err != nil {
		utility.ErrorHandler(err, false)
		return table, err
	}

	var rows []Row
	for _, model := range m {
		var columns []Column
		for _, definition := range cd {
			column := Column{}
			if definition.StaticValue != "" {
				column.Value = definition.StaticValue
			} else if definition.ComponentDefinition != nil {
				cModel := definition.ComponentDefinition.Struct
				for name, field := range definition.ComponentDefinition.Parameters {
					setValue := reflectModelParameterValue(model, field)
					cModel = setModelParameterValue(cModel, name, setValue)
				}
				column.Field = cModel
			} else if definition.ValueFromStructField != "" {
				column.Value = reflectModelParameterValue(model, definition.ValueFromStructField)
			} else {
				column.Value = ""
			}
			columns = append(columns, column)
		}
		rows = append(rows, Row{
			Columns: columns,
		})
	}

	var header []Column
	for _, definition := range cd {
		header = append(header, Column{Name: definition.Name})
	}
	return Table{
		Section:          section,
		Header:           header,
		Rows:             rows,
		Theme:            theme,
		Template:         template,
		PageNumber:       page,
		PageDisplayCount: numRows,
		T:                t,
	}, nil
}

func reflectModelParameterValue(model interface{}, parameterName string) string {
	rv := reflect.ValueOf(model)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return ""
	}
	// Dereference pointer
	rv = rv.Elem()
	// Lookup field by name
	fv := rv.FieldByName(parameterName)
	if !fv.IsValid() {
		return ""
	}
	// We expect a string field
	if fv.Kind() != reflect.String {
		// Convert interface to string
		return fmt.Sprintf("%v", fv.Interface())
	}
	return fv.String()
}

func setModelParameterValue(model page.Component, parameterName string, value string) page.Component {
	rv := reflect.ValueOf(model)
	if rv.Kind() != reflect.Ptr {
		return model
	}
	// Dereference pointer
	rv = rv.Elem()
	// Lookup field by name
	fv := rv.FieldByName(parameterName)
	if !fv.IsValid() {
		return model
	}
	// We expect a string field
	if fv.Kind() != reflect.String {
		// Convert interface to string
		return model
	}
	fv.SetString(value)
	return model
}
