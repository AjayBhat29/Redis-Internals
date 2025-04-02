package server

import (
	"log"
	"syscall"

	"github.com/AjayBhat29/Redis-Internals/config"
)

var concurrent_clients int = 0

func RunAsyncTCPServer() {
	log.Println("Starting asynchronous TCP server on ", config.HOST, ":", config.PORT)

	max_clients := 20000

	// Create EPOLL Event objects to hold events
	var events []syscall.EpollEvent = make([]syscall.EpollEvent, max_clients)

	// Create a new socket
}
