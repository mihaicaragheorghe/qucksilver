package server_test

import (
	"net"
	"testing"
	"time"

	"github.com/mihaicaragheorghe/qucksilver/internal/server"
)

func TestServer(t *testing.T) {
	addr := "localhost:13689"
	go func() {
		server.Start(addr)
	}()

	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Failed to connect to server %s:", err)
	}
	defer conn.Close()

	tests := []struct {
		input    string
		expected string
	}{
		{"+PING\r\n", "+PONG\r\n"},
		{"*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n", "+hey\r\n"},
		{"+FOO\r\n", "-ERR unknown command\r\n"},
	}

	for _, test := range tests {
		_, err := conn.Write([]byte(test.input))
		if err != nil {
			t.Fatalf("Failed to write command to server: %s", err)
		}

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatalf("Failed to read response from server :%s", err)
		}

		actual := string(buf[:n])
		if actual != test.expected {
			t.Errorf("Expected %s but received %s", test.expected, actual)
		}
	}
}
