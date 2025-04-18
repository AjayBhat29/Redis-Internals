package core

import (
	"errors"
	"io"
	"log"
	"strconv"
	"time"
)

var RESP_NIL []byte = []byte("$-1\r\n")

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

func evaluateSET(args []string, conn io.ReadWriter) error {
	if len(args) <= 1 {
		return errors.New("(error) ERR wrong number of arguments for 'set' command")
	}

	var key, value string
	var exDurationMs int64 = -1

	key, value = args[0], args[1]

	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "EX", "ex":
			i++
			if i == len(args) {
				return errors.New("(error) ERR syntax error")
			}

			exDurationSec, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return errors.New("(error) ERR value is not an integer or out of range")
			}
			exDurationMs = exDurationSec * 1000
		default:
			return errors.New("(error) ERR syntax error")
		}
	}

	// putting the k and value in a Hash Table
	Put(key, NewObj(value, exDurationMs))
	conn.Write([]byte("+OK\r\n"))
	return nil
}

func evaluateGET(args []string, conn io.ReadWriter) error {
	if len(args) != 1 {
		return errors.New("(error) ERR wrong number of arguments for 'get' command")
	}

	var key string = args[0]

	// Get the key from the hash table
	obj := Get(key)

	// if key does not exist, return RESP encoded nil
	if obj == nil {
		conn.Write(RESP_NIL)
		return nil
	}

	// if key already expired then return nil
	if obj.ExpiresAt != -1 && obj.ExpiresAt <= time.Now().UnixMilli() {
		conn.Write(RESP_NIL)
		return nil
	}

	// return the RESP encoded value
	conn.Write(Encode(obj.Value, false))
	return nil
}

func evaluateTTL(args []string, conn io.ReadWriter) error {
	if len(args) != 1 {
		return errors.New("(error) ERR wrong number of arguments for 'ttl' command")
	}

	var key string = args[0]

	obj := Get(key)

	// if key does not exist, return RESP encoded -2 denoting key does not exist
	if obj == nil {
		conn.Write([]byte(":-2\r\n"))
		return nil
	}

	// if object exist, but no expiration is set on it then send -1
	if obj.ExpiresAt == -1 {
		conn.Write([]byte(":-1\r\n"))
		return nil
	}

	// compute the time remaining for the key to expire and
	// return the RESP encoded form of it
	durationMs := obj.ExpiresAt - time.Now().UnixMilli()

	// if key expired, return -2
	if durationMs < 0 {
		conn.Write([]byte(":-2\r\n"))
		return nil
	}

	conn.Write(Encode(int64(durationMs/1000), false))
	return nil
}

func evaluateDEL(args []string, conn io.ReadWriter) error {
	var countDeleted int = 0

	for _, key := range args {
		if ok := Del(key); ok {
			countDeleted++
		}
	}

	conn.Write(Encode(countDeleted, false))
	return nil
}

func evaluateEXPIRE(args []string, conn io.ReadWriter) error {
	if len(args) <= 1 {
		return errors.New("(error) ERR wrong number of arguments for 'expire' command")
	}

	var key string = args[0]
	exDurationSec, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return errors.New("(error) ERR value is not an integer or out of range")
	}

	obj := Get(key)

	// 0 if the timeout was not set. e.g. key doesn't exist, or operation skipped due to the provided arguments
	if obj == nil {
		conn.Write([]byte(":0\r\n"))
		return nil
	}

	obj.ExpiresAt = time.Now().UnixMilli() + exDurationSec*1000

	// 1 if the timeout was set.
	conn.Write([]byte(":1\r\n"))
	return nil
}

func EvaluateAndRespond(cmd *RedisCmd, conn io.ReadWriter) error {
	log.Println("Command: ", cmd.Cmd)
	switch cmd.Cmd {
	case "PING":
		return evaluatePING(cmd.Args, conn)
	case "SET":
		return evaluateSET(cmd.Args, conn)
	case "GET":
		return evaluateGET(cmd.Args, conn)
	case "TTL":
		return evaluateTTL(cmd.Args, conn)
	case "DEL":
		return evaluateDEL(cmd.Args, conn)
	case "EXPIRE":
		return evaluateEXPIRE(cmd.Args, conn)
	default:
		return evaluatePING(cmd.Args, conn)
	}
}
