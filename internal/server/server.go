package server

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"

	"github.com/mihaicaragheorghe/qucksilver/internal/resp"
)

func Start(addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	log.Printf("Server started on %s\n", addr)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		data, err := read(conn)
		if err != nil {
			fmt.Printf("Error reading data: %s\n", err)
			return
		}

		cmd, args, err := resp.ProcessRESP(bufio.NewReader(bytes.NewReader(data)))
		if err != nil {
			log.Printf("Error processing RESP: %s\n", err)
			write(conn, fmt.Sprintf("-%s", "-ERR invalid RESP format\r\n"))
			return
		}
		fmt.Printf("Received command: %s\n", cmd)
		response := handleCommand(cmd, args)
		err = write(conn, response)
		if err != nil {
			fmt.Printf("Error writing data: %s", err)
			return
		}
	}
}

func handleCommand(cmd string, args []interface{}) string {
	switch cmd {
	case "PING":
		return "+PONG\r\n"
	case "ECHO":
		if len(args) > 0 {
			return fmt.Sprintf("+%s\r\n", args[0])
		}
		return "-ERR wrong number of arguments for 'echo' command\r\n"
	default:
		return "-ERR unknown command\r\n"
	}
}

func read(conn net.Conn) ([]byte, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	log.Printf("Received %d bytes\n", n)
	return buf[:n], nil
}

func write(conn net.Conn, s string) error {
	n, err := conn.Write([]byte(s))
	if err != nil {
		return err
	}
	log.Printf("Sent %d bytes\n", n)
	return nil
}
