package main

import (
	"flag"
	"log"

	"github.com/AjayBhat29/Redis-Internals/config"
	"github.com/AjayBhat29/Redis-Internals/server"
)

func setupFlags() {
	flag.StringVar(&config.HOST, "host", "0.0.0.0", "host for the server")
	flag.IntVar(&config.PORT, "port", 8379, "port for the server")
	flag.Parse()
}

func main() {
	setupFlags()
	log.Println("Starting server on", config.HOST, ":", config.PORT)
	server.RunSyncTCPServer()
}
