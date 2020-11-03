package routes

import (
	"fmt"
	"github.com/markustenghamn/beubo/pkg/structs"
	"github.com/markustenghamn/beubo/pkg/structs/page"
	"github.com/markustenghamn/beubo/pkg/structs/page/component"
	"github.com/markustenghamn/beubo/pkg/utility"
	"html/template"
	"net/http"
)

// Admin is the default admin route and template
func (br *BeuboRouter) Admin(w http.ResponseWriter, r *http.Request) {
	table, err := component.MakeTable(br.DB, []structs.Site{}, []component.ColumnDefinition{
		{Name: "ID", ValueFromStructField: "ID"},
		{Name: "Site", ValueFromStructField: "Title"},
		{Name: "Domain", ValueFromStructField: "Domain"},
		{Name: "", ComponentDefinition: &page.ComponentDefinition{
			Parameters: map[string]page.ComponentParameterDefinition{
				"Link": page.ComponentParameterDefinition{
					StaticValue: "",
					StructField: "",
					// TODO maybe a computed field is needed?
				},
			},
			// TODO fix schema here
			//Link:    template.URL(fmt.Sprintf("%s://%s/", "http", site.Domain)),
			//Class:   "button-primary",
			//Content: "View",
			//T:       br.Renderer.T,
		}},
		//{Name: "", Field: component.Button{
		//	Link:    template.URL(fmt.Sprintf("/admin/sites/a/%s", sid)),
		//	Class:   "button-primary",
		//	Content: "Manage",
		//	T:       br.Renderer.T,
		//}},
		//{Name: "", Field: component.Button{
		//	Link:    template.URL(fmt.Sprintf("/admin/sites/edit/%s", sid)),
		//	Class:   "button-primary",
		//	Content: "Edit",
		//	T:       br.Renderer.T,
		//}},
		//{Name: "", Field: component.Button{
		//	Link:    template.URL(fmt.Sprintf("/admin/sites/delete/%s", sid)),
		//	Class:   "button-clear",
		//	Content: "Delete",
		//	T:       br.Renderer.T,
		//}}
	}, 10, 0, "main", "", "", br.Renderer.T)

	utility.ErrorHandler(err, false)

	pageData := structs.PageData{
		Template: "admin.page",
		Title:    "Admin - Sites",
		Components: []page.Component{
			component.Button{
				Section: "main",
				Link:    template.URL("/admin/sites/add"),
				Class:   "button-primary",
				Content: "Add",
				T:       br.Renderer.T,
			},
			table,
		},
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

func (br *BeuboRouter) Settings(w http.ResponseWriter, r *http.Request) {
	var settings []structs.Setting

	if err := br.DB.Find(&settings).Error; err != nil {
		utility.ErrorHandler(err, false)
	}

	var rows []component.Row
	for _, setting := range settings {
		sid := fmt.Sprintf("%d", setting.ID)
		rows = append(rows, component.Row{
			Columns: []component.Column{
				{Name: "ID", Value: sid},
				{Name: "Key", Value: setting.Key},
				{Name: "Value", Value: setting.Value},
			},
		})
	}

	table := component.Table{
		Section: "main",
		Header: []component.Column{
			{Name: "ID"},
			{Name: "Key"},
			{Name: "Value"},
		},
		Rows:             rows,
		PageNumber:       1,
		PageDisplayCount: 10,
		T:                br.Renderer.T,
	}

	pageData := structs.PageData{
		Template: "admin.page",
		Title:    "Admin - Settings",
		Components: []page.Component{
			table,
		},
	}
	br.Renderer.RenderHTMLPage(w, r, pageData)
}

func (br *BeuboRouter) Users(w http.ResponseWriter, r *http.Request) {
	var users []structs.User

	if err := br.DB.Find(&users).Error; err != nil {
		utility.ErrorHandler(err, false)
	}

	var rows []component.Row
	for _, user := range users {
		sid := fmt.Sprintf("%d", user.ID)
		rows = append(rows, component.Row{
			Columns: []component.Column{
				{Name: "ID", Value: sid},
				{Name: "Email", Value: user.Email},
			},
		})
	}

	table := component.Table{
		Section: "main",
		Header: []component.Column{
			{Name: "ID"},
			{Name: "Email"},
		},
		Rows:             rows,
		PageNumber:       1,
		PageDisplayCount: 10,
		T:                br.Renderer.T,
	}

	pageData := structs.PageData{
		Template: "admin.page",
		Title:    "Admin - Users",
		Components: []page.Component{
			table,
		},
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)

}

func (br *BeuboRouter) GetPlugins(w http.ResponseWriter, r *http.Request) {
	var rows []component.Row
	for plugin := range *br.Plugins {
		rows = append(rows, component.Row{
			Columns: []component.Column{
				{Name: "Name", Value: plugin},
			},
		})
	}

	table := component.Table{
		Section: "main",
		Header: []component.Column{
			{Name: "Name"},
		},
		Rows:             rows,
		PageNumber:       1,
		PageDisplayCount: 10,
		T:                br.Renderer.T,
	}

	pageData := structs.PageData{
		Template: "admin.page",
		Title:    "Admin - Plugins",
		Components: []page.Component{
			table,
		},
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}
