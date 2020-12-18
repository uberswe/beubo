package plugin

import "plugin"

// Handler holds all the needed information for plugins to function in Beubo
type Handler struct {
	Plugins map[string]Plugin
}

// Plugin represents data for a single plugin in Beubo
type Plugin struct {
	Plugin     *plugin.Plugin
	Definition string
	Data       map[string]string
}
