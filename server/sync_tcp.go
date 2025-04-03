package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/AjayBhat29/Redis-Internals/config"
	"github.com/AjayBhat29/Redis-Internals/core"
)

func readCommand(conn io.ReadWriter) (*core.RedisCmd, error) {
	var buf []byte = make([]byte, 512)
	n, err := conn.Read(buf[:])
	if err != nil {
		return nil, err
	}

	tokens, err := core.DecodeArrayString(buf[:n])
	if err != nil {
		return nil, err
	}

	return &core.RedisCmd{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil
}

func respondError(err error, conn io.ReadWriter) {
	conn.Write([]byte(fmt.Sprintf("-%s\r\n", err.Error())))
}

func respond(cmd *core.RedisCmd, conn io.ReadWriter) {
	err := core.EvaluateAndRespond(cmd, conn)
	if err != nil {
		respondError(err, conn)
	}
}

func RunSyncTCPServer() {
	log.Println("Starting synchronous TCP server on", config.HOST, ":", config.PORT)

	var concurrent_clients int = 0

	listener, err := net.Listen("tcp", config.HOST+":"+strconv.Itoa(config.PORT))
	if err != nil {
		log.Println("Error starting server:", err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error: ", err)
		}

		concurrent_clients++
		log.Println("Client connected with address", conn.RemoteAddr().String(), "Number of clients:", concurrent_clients)

		for {
			cmd, err := readCommand(conn)
			if err != nil {
				conn.Close()
				concurrent_clients--
				log.Println("Client disconnected with address", conn.RemoteAddr().String(), "Number of clients:", concurrent_clients)
				if err == io.EOF {
					break
				}
				log.Println("Error reading command. Got:", err)
			}

			respond(cmd, conn)
		}
	}

}
