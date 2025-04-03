package server

import (
	"log"
	"net"
	"syscall"

	"github.com/AjayBhat29/Redis-Internals/config"
	"github.com/AjayBhat29/Redis-Internals/core"
)

var concurrent_clients int = 0

func RunAsyncTCPServer() error {
	log.Println("Starting asynchronous TCP server on ", config.HOST, ":", config.PORT)

	max_clients := 20000

	// Create EPOLL Event objects to hold events
	var events []syscall.EpollEvent = make([]syscall.EpollEvent, max_clients)

	// Create a new socket
	serverFD, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(serverFD) //Once the function execution is done, close the socket

	//Set the socket to operate in a non-blocking mode
	if err = syscall.SetNonblock(serverFD, true); err != nil {
		return err
	}

	//Bind IP and PORT
	ip4 := net.ParseIP(config.HOST)
	if err = syscall.Bind(serverFD, &syscall.SockaddrInet4{
		Port: config.PORT,
		Addr: [4]byte{ip4[0], ip4[1], ip4[2], ip4[3]},
	}); err != nil {
		return err
	}

	//Start listning
	if err = syscall.Listen(serverFD, max_clients); err != nil {
		return err
	}

	//Async IO starts here

	//Create an epoll instance
	epollFD, err := syscall.EpollCreate1(0)
	if err != nil {
		log.Fatal(err)
	}
	defer syscall.Close(epollFD)

	// Specify the events we want to get hints about and set the socket to be monitored
	var socketServerEvent syscall.EpollEvent = syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(serverFD),
	}

	// Registering in the Epoll, add this file descriptor to be monitored for incoming connections
	if err = syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, serverFD, &socketServerEvent); err != nil {
		return err
	}

	// Start the event loop
	for {
		//See if any FD is ready for an IO
		nevents, e := syscall.EpollWait(epollFD, events[:], -1)
		if e != nil {
			continue
		}

		for i := 0; i < nevents; i++ {
			// if the server is ready for an IO
			if int(events[i].Fd) == serverFD {
				//accept incoming connection from client
				fd, _, err := syscall.Accept(serverFD)
				if err != nil {
					log.Println("Error accepting connection:", err)
					continue
				}

				concurrent_clients++
				syscall.SetNonblock(fd, true)

				// Add this new TCP connection to be monitored by the epoll instance
				var socketClientEvent syscall.EpollEvent = syscall.EpollEvent{
					Events: syscall.EPOLLIN,
					Fd:     int32(fd),
				}
				if err := syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, fd, &socketClientEvent); err != nil {
					log.Fatal(err)
				}
			} else {
				comm := core.FDComm{Fd: int(events[i].Fd)}
				cmd, err := readCommand(comm)
				if err != nil {
					syscall.Close(int(events[i].Fd))
					concurrent_clients--
					continue
				}
				respond(cmd, comm)
			}
		}
	}
}
