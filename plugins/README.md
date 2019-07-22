## Plugins

Beubo supports plugins and it should use the Go plugins package https://golang.org/pkg/plugin/

More documentation, sample code and a proper interface will be added.

## Build a plugin

To run a plugin it needs to be compiled, to build the plugin run:

```
go build -buildmode=plugin
```