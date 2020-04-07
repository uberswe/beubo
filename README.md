## Beubo

Simply run `go run cmd/beubo/main.go` to get started.

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
