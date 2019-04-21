package main

import "github.com/markustenghamn/beubo"

func main() {
	beubo.Init()
	beubo.Seed()
	beubo.InitRoutes()
}
