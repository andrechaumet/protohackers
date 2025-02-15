package main

import (
	"errors"
	"io"
	"log"
	"net"
	"time"
)

// protohackers.com/problem/0
func main() {
	ln := setup(":8080")
	defer ln.Close()
	handleConns(ln)
}

func setup(address string) net.Listener {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	return ln
}

func handleConns(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(10 * time.Second))
	buf := make([]byte, 1024)
	if n, err := readConn(conn, buf); err != nil {
		log.Println("Error reading:", err)
	} else if err := writeConn(conn, buf[:n]); err != nil {
		log.Println("Error writing:", err)
	}
}

func readConn(conn net.Conn, buf []byte) (int, error) {
	n, err := conn.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		return n, err
	}
	return n, nil
}

func writeConn(conn net.Conn, data []byte) error {
	_, err := conn.Write(data)
	return err
}
