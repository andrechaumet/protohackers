package main

import (
	"log"
	"net"
	"strings"
	"sync"
)

type database struct {
	lock sync.Mutex
	data map[string]string
}

const v = "omg hii andy was here"

// protohackers.com/problem/4
func main() {
	addr := resolve(":8080")
	conn := start(addr)
	defer conn.Close()
	db := database{sync.Mutex{}, make(map[string]string)}
	listen(conn, &db)
}

func resolve(port string) *net.UDPAddr {
	addr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		panic(err)
	}
	return addr
}

func start(addr *net.UDPAddr) *net.UDPConn {
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	return conn
}

func listen(conn *net.UDPConn, db *database) {
	for {
		buf := make([]byte, 1000)
		_, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		go handle(string(buf), db)
	}
}

func handle(request string, db *database) string {
	if is, val := version(request); is {
		return val
	}
	if is, val := insert(request, db); is {
		log.Println(val)
		return val
	} else {
		return db.data[request]
	}
}

func version(data string) (bool, string) {
	if strings.Contains(data, "version") {
		return true, v
	}
	return false, ""
}

func insert(data string, d *database) (bool, string) {
	for i := 0; i < len(data); i++ {
		if data[i] == '=' {
			key := data[:i]
			value := data[1+i : len(data)]
			d.data[key] = value
			return true, d.data[key]
		}
	}
	return false, ""
}
