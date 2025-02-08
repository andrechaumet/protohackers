package main

import (
	"bytes"
	"net"
)

func main() {
	addr := resolve(":8082")
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
	return conn
}

func listen(conn *net.UDPConn) {
	for {
		buf := make([]byte, 1024)
		if _, addr, err := conn.ReadFromUDP(buf); err != nil {
			continue
		} else {
			request := string(bytes.Trim(buf, "\x00"))
			go handle(request, addr)
		}
	}
}

func handle(request string, addr *net.UDPAddr) {
	if parts, err := Validate(request); err != nil {
		return
	} else {
		Types[parts[0]](parts)
	}
}
