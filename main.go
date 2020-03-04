package main

import (
	"github.com/duyk16/notifying-server/config"
	"github.com/duyk16/notifying-server/server"
	"github.com/duyk16/notifying-server/storage"
)

func main() {
	config.Init()
	storage.Init()
	server.Init()
}
