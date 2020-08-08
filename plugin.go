package beubo

import (
	"log"
	"path/filepath"
	"plugin"
)

// This is the plugin handler for Beubo
// The initial function will always be Register
// We need a fairly dynamic way to hook any Beubo function into the plugins
// Most functions can be async but some will cause Beubo to wait for them to finish.
// Analytics does not need Beubo to wait as it doesn't affect the flow
// An ecommerce plugin would need Beubo to wait since it affects what is displayed on the page

func loadPlugins() {
	// The plugins (the *.so files) must be in a 'plugins' sub-directory
	allPlugins, err := filepath.Glob("plugins/*.so")
	if err != nil {
		panic(err)
	}

	for _, filename := range allPlugins {
		p, err := plugin.Open(filename)
		if err != nil {
			log.Println(err)
			continue
		}

		symbol, err := p.Lookup("Register")
		if err != nil {
			panic(err)
		}

		registerFunc, ok := symbol.(func() map[string]string)
		if !ok {
			panic("Plugin has no 'Register() map[string]string' function")
		}

		pluginData := registerFunc()
		log.Println("Loading plugins")
		log.Println(filename, pluginData)
	}
}
