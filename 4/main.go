package main

import (
	"log"
	"net"
	"sync"
)

// protohackers.com/problem/4
type database struct {
	lock sync.Mutex
	data map[string]string
}

func (db *database) save(key string, value string) string {
	db.lock.Lock()
	defer db.lock.Unlock()
	db.data[key] = value
	return db.data[key]
}

var db = database{sync.Mutex{}, make(map[string]string)}

func main() {
	addr := resolve(":8080")
	conn := start(addr)
	defer conn.Close()
	listen(conn)
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
	db.save("version", "v.YOLO")
	return conn
}

func listen(conn *net.UDPConn) {
	for {
		buf := make([]byte, 1001)
		_, addr, err := conn.ReadFromUDP(buf)
		if err != nil || len(buf) > 1000 {
			continue
		}
		go handle(string(buf), addr, conn)
	}
}

func handle(request string, addr *net.UDPAddr, conn *net.UDPConn) {
	response := process(request)
	if _, err := conn.WriteToUDP([]byte(response), addr); err != nil {
		log.Println("Error while returning request data to OP")
	}
}

func process(request string) string {
	if is, val := insert(request); is {
		return val
	} else {
		return db.data[request]
	}
}

func insert(request string) (bool, string) {
	for i := 0; i < len(request); i++ {
		if request[i] == '=' {
			key, value := request[:i], request[1+i:]
			return true, db.save(key, value)
		}
	}
	return false, ""
}
