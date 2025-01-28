package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

// protohackers.com/problem/1
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
	buf := make([]byte, 1024)
	for {
		if n, err := readConn(conn, buf); err != nil {
			log.Println("Error reading:", err)
			return
		} else {
			err := writeConn(buf, n)
			if err != nil {
				log.Println("Error writing:", err)
				return
			}
		}
	}
}

func readConn(conn net.Conn, buf []byte) (int, error) {
	n, err := conn.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		return n, err
	}
	return n, nil
}

func writeConn(buf []byte, n int) error {
	if req, err := extract(string(buf[:n])); err != nil {
		log.Println(err)
		return fmt.Errorf("could not extract data")
	} else {
		eval := prime(req)
		if err := writeConn(nil, 0); err != nil {
			return err
		} else if !eval {
			return fmt.Errorf("%.2f is not prime", req)
		}
	}
	return nil
}

func extract(input string) (float64, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		return 0, fmt.Errorf("invalid JSON: %v", err)
	}

	if method, ok := data["method"].(string); !ok || method != "isPrime" {
		return 0, fmt.Errorf("no key 'method' 'isPrime' present")
	}

	if number, ok := data["number"].(float64); ok {
		return float64(number), nil
	}
	return 0, fmt.Errorf("key 'number' is not numeric")
}

func prime(ev float64) bool {
	if ev < 2 || ev != float64(int(ev)) {
		return false
	}
	n := int(ev)
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}
