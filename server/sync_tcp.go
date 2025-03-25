package server

import (
	"io"
	"log"
	"net"
	"strconv"

	"github.com/AjayBhat29/Redis-Internals/config"
)

func readCommand(c net.Conn) (string, error) {
	var buf []byte = make([]byte, 512)
	n, err := c.Read(buf[:])
	if err != nil {
		return "", err
	}
	return string(buf[0:n]), nil
}

func respond(cmd string, c net.Conn) error {
	if _, err := c.Write([]byte(cmd)); err != nil {
		return err
	}
	return nil
}

func RunSyncTCPServer() {
	log.Println("Starting synchronous TCP server on", config.HOST, ":", config.PORT)

	var concurrent_clients int = 0

	listener, err := net.Listen("tcp", config.HOST+":"+strconv.Itoa(config.PORT))
	if err != nil {
		panic(err)
	}

	for {
		c, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		concurrent_clients++
		log.Println("Client connected with address", c.RemoteAddr().String(), "Number of clients:", concurrent_clients)

		for {
			cmd, err := readCommand(c)
			if err != nil {
				c.Close()
				concurrent_clients--
				log.Println("Client disconnected with address", c.RemoteAddr().String(), "Number of clients:", concurrent_clients)
				if err == io.EOF {
					break
				}
				log.Println("Error reading command. Got:", err)
			}

			log.Println("Command received:", cmd)
			if err = respond(cmd, c); err != nil {
				log.Print("Error responding. Got:", err)
			}
		}
	}

}
