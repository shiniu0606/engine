package main

import (
	"flag"
	server "github.com/shiniu0606/engine/server/servercenter"
)

func main() {
	flag.Parse()

	server.CreateDBTable()
}