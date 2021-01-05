package plugin

import (
	"github.com/uberswe/beubo/pkg/structs"
	"log"
	"path/filepath"
	"plugin"
)

// Load is a function that loads any available plugins in the plugins/ folder
func Load(pluginHandler Handler) Handler {
	// The plugins (the *.so files) must be in a 'plugins' sub-directory
	allPlugins, err := filepath.Glob("plugins/*.so")
	if err != nil {
		panic(err)
	}

	log.Println("Loading plugins...")
	for _, filename := range allPlugins {
		plug, err := plugin.Open(filename)
		if err != nil {
			log.Println(err)
			continue
		}

		symbol, err := plug.Lookup("Register")
		if err != nil {
			log.Println("Plugin has no 'Register() map[string]string' function")
			continue
		}

		registerFunc, ok := symbol.(func() map[string]string)
		if !ok {
			log.Println("Plugin has no 'Register() map[string]string' function")
			continue
		}

		pluginData := registerFunc()
		log.Println(filename, pluginData)

		if pluginHandler.Plugins == nil {
			pluginHandler.Plugins = map[string]Plugin{}
		}

		pluginHandler.Plugins[filename] = Plugin{
			Plugin:     plug,
			Definition: pluginData["identifier"],
			Data:       pluginData,
		}
	}

	// TODO the following code requires an active database connection but what if a plugin wants to modify a connection? It would be good if this could be moved so the loading of plugins happens before the active check is performed.
	if pluginHandler.Plugins != nil {
		sites := structs.FetchSites(pluginHandler.DB)
		var plugins []string
		var existingSites []uint
		for _, s := range sites {
			existingSites = append(existingSites, s.ID)
		}
		for _, p := range pluginHandler.Plugins {
			plugins = append(plugins, p.Definition)
			for _, s := range sites {
				ps := PluginSite{}
				pluginHandler.DB.Where("plugin_identifier = ?", p.Definition).Find(&ps)
				if ps.ID <= 0 {
					ps.Active = false
					ps.SiteID = s.ID
					ps.PluginIdentifier = p.Definition
					pluginHandler.DB.Create(&p)
				}
			}
		}
		if len(plugins) > 0 {
			pluginHandler.DB.Where("plugin_identifier NOT IN ?", plugins).Delete([]PluginSite{})
		} else {
			// delete everything if there are no plugins
			pluginHandler.DB.Where("1=1").Delete([]PluginSite{})
		}
		if len(existingSites) > 0 {
			pluginHandler.DB.Where("site_id NOT IN ?", existingSites).Delete([]PluginSite{})
		} else {
			// delete everything if there are no sites
			pluginHandler.DB.Where("1=1").Delete([]PluginSite{})
		}
	}

	return pluginHandler
}
