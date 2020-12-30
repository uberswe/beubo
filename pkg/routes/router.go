package routes

import (
	"github.com/uberswe/beubo/pkg/plugin"
	"github.com/uberswe/beubo/pkg/template"
	"gorm.io/gorm"
)

// BeuboRouter holds parameters relevant to the router
type BeuboRouter struct {
	DB            *gorm.DB
	Renderer      *template.BeuboTemplateRenderer
	PluginHandler *plugin.Handler
}
