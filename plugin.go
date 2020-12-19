package beubo

import (
	beuboplugin "github.com/uberswe/beubo/pkg/plugin"
	"log"
	"path/filepath"
	"plugin"
)

// This is the plugin handler for Beubo
// The initial function will always be Register
// We need a fairly dynamic way to hook any Beubo function into the plugins
// Most functions can be async but some will cause Beubo to wait for them to finish.
// For example, an analytics plugin does not need Beubo to wait as it doesn't affect the flow
// But an ecommerce plugin would need Beubo to wait since it affects what is displayed on the page

var pluginHandler beuboplugin.Handler

func loadPlugins() {
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
			pluginHandler.Plugins = map[string]beuboplugin.Plugin{}
		}

		pluginHandler.Plugins[filename] = beuboplugin.Plugin{
			Plugin: plug,
			// TODO check for this key
			Definition: pluginData["identifier"],
			Data:       pluginData,
		}
	}
}
