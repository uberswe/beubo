package plugin

import "plugin"

type Handler struct {
	Plugins map[string]Plugin
}

type Plugin struct {
	Plugin     *plugin.Plugin
	Definition string
	Data       map[string]string
}
