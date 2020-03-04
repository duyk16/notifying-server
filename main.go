package main

import (
	"github.com/duyk16/notifying-server/config"
	"github.com/duyk16/notifying-server/server"
)

func main() {
	config.Init()
	server.Init()
}
