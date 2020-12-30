package plugin

import (
	"github.com/uberswe/beubo/pkg/structs"
	"gorm.io/gorm"
	"plugin"
)

// Handler holds all the needed information for plugins to function in Beubo
type Handler struct {
	DB      *gorm.DB
	Plugins map[string]Plugin
}

// Plugin represents data for a single plugin in Beubo
type Plugin struct {
	Plugin     *plugin.Plugin
	Definition string
	Data       map[string]string
}

type PluginSite struct {
	gorm.Model
	Site             structs.Site
	SiteID           uint
	PluginIdentifier string
	Active           bool
}

func FetchPluginSites(db *gorm.DB, plugin string) (ps []PluginSite) {
	db.Preload("Site").Where("plugin_identifier = ?", plugin).Find(&ps)
	return ps
}
