## Plugins

Beubo supports plugins and it should use the Go plugins package https://golang.org/pkg/plugin/

More documentation, sample code and a proper interface will be added.

## Build a plugin

To run a plugin it needs to be compiled, to build the plugin run:

```
go build -buildmode=plugin
```

## Available functions

Currently the only function available is the `Register` function which is called when Beubo starts.
This can be used to tell Beubo about the plugin.

## Planned functions

 - Page request
 - Template
 - Save page
 - New page
 - New site
 - Save site
 - Login
 - Logout
 - Install started
 - Install finished
 - Beubo is running, a run function that starts running if Beubo is active
 - Shutdown