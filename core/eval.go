package core

import (
	"errors"
	"io"
	"log"
)

func evaluatePING(args []string, conn io.ReadWriter) error {
	var b []byte

	if len(args) >= 2 {
		return errors.New("ERROR: wrong number of arguments for 'PING' command")
	}

	if len(args) == 0 {
		b = Encode("PONG", true)
	} else {
		b = Encode(args[0], false)
	}

	_, err := conn.Write(b)
	return err
}

func EvaluateAndRespond(cmd *RedisCmd, conn io.ReadWriter) error {
	log.Println("Command: ", cmd.Cmd)
	switch cmd.Cmd {
	case "PING":
		return evaluatePING(cmd.Args, conn)
	default:
		return errors.New("ERROR: unknown command")
	}
}
