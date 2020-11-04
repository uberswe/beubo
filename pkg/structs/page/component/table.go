package component

import (
	"fmt"
	pluralize "github.com/gertd/go-pluralize"
	"github.com/jinzhu/gorm"
	"github.com/markustenghamn/beubo/pkg/structs/page"
	"github.com/markustenghamn/beubo/pkg/utility"
	"html/template"
	"log"
	"reflect"
	"strconv"
	"strings"
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
	mType := reflect.TypeOf(m)
	results := reflect.New(reflect.SliceOf(mType)).Interface()
	table := Table{}
	var rows []Row
	if err := db.Model(m).Limit(numRows).Offset(page).Find(results).Error; err != nil {
		utility.ErrorHandler(err, false)
		return table, err
	}
	switch reflect.TypeOf(results).Elem().Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(results).Elem()
		for i := 0; i < s.Len(); i++ {
			model := s.Index(i)
			var columns []Column
			for _, definition := range cd {
				column := Column{}
				if definition.StaticValue != "" {
					column.Value = definition.StaticValue
				} else if definition.ComponentDefinition != nil {
					cModel := definition.ComponentDefinition.Struct
					for name, field := range definition.ComponentDefinition.Parameters {
						setValue := ""
						if field.StaticValue != "" {
							setValue = field.StaticValue
						} else {
							setValue = reflectModelParameterValue(model, field.StructField)
						}
						if field.ComputedField != nil {
							setValue = field.ComputedField(setValue)
						}
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
	log.Printf("%v\n", model)
	rv := reflect.ValueOf(model)
	fields := reflect.TypeOf(model)
	values := reflect.ValueOf(model)
	num := fields.NumField()

	for i := 0; i < num; i++ {
		field := fields.Field(i)
		value := values.Field(i)
		fmt.Print("Type:", field.Type, ",", field.Name, "=", value, "\n")

		switch value.Elem().Kind() {
		case reflect.String:
			v := value.Elem().String()
			fmt.Print(v, "\n")
		case reflect.Int:
			v := strconv.FormatInt(value.Elem().Int(), 10)
			fmt.Print(v, "\n")
		case reflect.Int32:
			v := strconv.FormatInt(value.Elem().Int(), 10)
			fmt.Print(v, "\n")
		case reflect.Int64:
			v := strconv.FormatInt(value.Elem().Int(), 10)
			fmt.Print(v, "\n")
		default:
			fmt.Printf("Not support type of struct %s", value.Kind())
		}
	}

	log.Println(values)

	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		// Lookup field by name
		fv := rv.FieldByName(parameterName)
		if !fv.IsValid() {
			log.Println(rv.Kind())
			log.Println(parameterName)
			log.Println(fv.IsValid())
			return ""
		}
		// We expect a string field
		if fv.Kind() != reflect.String {
			// Convert interface to string
			return fmt.Sprintf("%v", fv.Interface())
		}
		return fv.String()
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

func getTable(m interface{}) string {
	p := pluralize.NewClient()
	t := reflect.TypeOf(m)
	n := t.Name()
	if t.Kind() == reflect.Ptr {
		n = t.Elem().Name()
	}
	n = strings.ToLower(n)
	if !p.IsPlural(n) {
		n = p.Plural(n)
	}
	return n
}
