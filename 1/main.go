package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

// protohackers.com/problem/1
func main() {
	ln := setup(":8080")
	defer ln.Close()
	listen(ln)
}

func setup(address string) net.Listener {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	return ln
}

func listen(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		if _, err := read(conn, buf); err != nil {
			log.Println("Error reading:", err)
			return
		}
		val, err := extract(string(buf))
		if err != nil {
			log.Println("Error extracting:", err)
			return
		}
		eval := prime(val)
		if err := write(conn, eval); err != nil {
			return
		}
	}
}

func read(conn net.Conn, buf []byte) (int, error) {
	n, err := conn.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		return n, err
	}
	return n, nil
}

func write(conn net.Conn, ans bool) error {
	res := fmt.Sprintf(`{"method": "isPrime","number": %t}`, ans)
	_, err := conn.Write([]byte(res))
	return err
}

func extract(input string) (float64, error) {
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "\x00", "")
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		return 0, fmt.Errorf("invalid JSON: %v", err)
	}

	if method, ok := data["method"].(string); !ok || method != "isPrime" {
		return 0, fmt.Errorf("no key 'isPrime' present")
	}

	if number, ok := data["number"].(float64); ok {
		return float64(number), nil
	}
	return 0, fmt.Errorf("value 'number' is not numeric")
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
