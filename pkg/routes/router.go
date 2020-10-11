package routes

import (
	"github.com/jinzhu/gorm"
	"github.com/markustenghamn/beubo/pkg/template"
)

type BeuboRouter struct {
	DB       *gorm.DB
	Renderer *template.BeuboTemplateRenderer
	Plugins  *map[string]map[string]string
}
