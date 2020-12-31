package beubo

import (
	beuboplugin "github.com/uberswe/beubo/pkg/plugin"
)

// This is the plugin handler for Beubo
// The initial function will always be Register
// We need a fairly dynamic way to hook any Beubo function into the plugins
// Most functions can be async but some will cause Beubo to wait for them to finish.
// For example, an analytics plugin does not need Beubo to wait as it doesn't affect the flow
// But an ecommerce plugin would need Beubo to wait since it affects what is displayed on the page

var pluginHandler beuboplugin.Handler

func loadPlugins() {
	pluginHandler = beuboplugin.Handler{
		DB: DB,
	}
	pluginHandler = beuboplugin.Load(pluginHandler)
}
