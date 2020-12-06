package routes

import (
	"github.com/jinzhu/gorm"
	"github.com/uberswe/beubo/pkg/template"
)

// BeuboRouter holds parameters relevant to the router
type BeuboRouter struct {
	DB       *gorm.DB
	Renderer *template.BeuboTemplateRenderer
	Plugins  *map[string]map[string]string
}
