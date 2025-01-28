package main

import (
	"errors"
	"io"
	"log"
	"net"
	"sync"
)

// protohackers.com/problem/2

/*
	byte:  |  0   |  1     2     3     4  |  5     6     7     8  |
	type:  | char |         int32         |         int32         |
*/

const size = 9

var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, size)
	},
}

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
	buf := bufferPool.Get().([]byte)
	data, err := readConn(conn, buf)
	if err != nil || data != size {
		return
	}
	process(buf)
}

func process(data []byte) error {
	operation := rune(data[0])
	if operation != 'I' && operation != 'Q' {

	}
}

/*
49
time := extract(data[0:4])0
price := extract(data[5:9])
--
51
minTime := extract(data[0:4])
maxTime := extract(data[5:9])
--
    Hexadecimal:                 Decoded:
<-- 49 00 00 30 39 00 00 00 65   I 12345 101
<-- 49 00 00 30 3a 00 00 00 66   I 12346 102
<-- 49 00 00 30 3b 00 00 00 64   I 12347 100
<-- 49 00 00 a0 00 00 00 00 05   I 40960 5
<-- 51 00 00 30 00 00 00 40 00   Q 12288 16384
--> 00 00 00 65                  101
*/

func extract(data []byte) int32 {
	return int32(bytesArray[1])<<24 | int32(bytesArray[2])<<16 | int32(bytesArray[3])<<8 | int32(bytesArray[4])
}

func readConn(conn net.Conn, buf []byte) (int, error) {
	n, err := conn.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		return n, err
	}
	return n, nil
}
