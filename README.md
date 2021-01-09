# Beubo

**Beubo is in Alpha and not recommended for production use, expect breaking changes and bugs**

I created Beubo to get better at Go. To learn more and to make it easier to get going with
new projects. None of the platforms or libraries in the Go ecosystem felt right for me.
That's why I set out to make my own CMS/Library, Beubo.

Beubo is a CMS that aims to be easy to use and written in Go. I wanted it 
to be as easy to use as Wordpress but with much better performance and with support
for multiple websites right from the start. I try to keep the capabilities of Beubo 
as small as possible. I hope I can make Beubo easy to build on using plugins so that 
it can be used for anything and everything.

Here are a few of the features I want to support:
 - Site management, routing based on domain
 - Page creation, editing, deletion
 - Themes
 - Plugins
 - User management with roles and permissions
 
That's pretty much it.

Simply run `go run cmd/beubo/main.go` to get started.

## Database

Beubo uses [GORM](https://gorm.io/) to handle database operations. Currently I am supporting sqlite3 and mysql but other drivers may work but have not been tested.

## Installation

When running Beubo for the first time an installation page will open at the specified port. The 
page asks for various details needed to configure your site including database details. You will 
need to create a database on a MariaDB server and provide details
so that Beubo can connect to it. You can also use sqlite3 but I only recommend it for local development.

Once the installation is complete it will no longer be available, delete the .env file to redo the 
installation process. To start with a fresh database simply truncate your current database and it will 
auto migrate and seed a fresh database.

## CLI options

```
-port=8080      Allows you to specify which port Beubo should listen on
```

## Templating

Beubo uses the go html templates to build pages. These templates use functions to render sections of 
content which plugins can hook into when a request is made.

## Plugins

Beubo supports go plugins. Simply place your `.so` under `/plugins` and Beubo will try to load this
plugin as it starts. A plugin will need to expose a `Register` method in order to run. Please see
the example plugin to learn more.