package server

import (
	"bufio"
	"fmt"
	"net"

	"github.com/BenDowswell/go-kv-store/kvstore"
	"github.com/BenDowswell/go-kv-store/protocol"
)

// ServeTCP binds to addr and serves the plain-text KV protocol over TCP.
// Each connection is handled in its own goroutine.
func ServeTCP(s kvstore.Store, addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return ServeTCPListener(s, ln)
}

// ServeTCPListener serves on an existing listener.
// Separating bind from serve lets tests pass in a :0 listener and read the
// actual port without a race between bind and the goroutine starting.
func ServeTCPListener(s kvstore.Store, ln net.Listener) error {
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go handleConn(conn, s)
	}
}

func handleConn(conn net.Conn, s kvstore.Store) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		cmd, err := protocol.Parse(scanner.Text())
		if err != nil {
			fmt.Fprintf(conn, "ERR %s\n", err)
			continue
		}
		fmt.Fprintf(conn, "%s\n", execute(s, cmd))
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(conn, "ERR read: %s\n", err)
	}
}
