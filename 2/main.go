package main

import (
	"errors"
	"git
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

type insertion struct {
	time   int32
	amount int32
}

type selection struct {
	start int32
	end   int32
}

func newInsert(data []byte) {

}

func newSelect(data []byte) {

}

const size = 9

var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, size)
	},
}

/*
func NewDateStore() *redblacktree.Tree {
	return &DateStore{
		tree: redblacktree.New(),
	}
}
*/

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
	nodes := redblacktree.Tree{}
	for {
		buf := bufferPool.Get().([]byte)
		data, err := read(conn, buf)
		if err != nil || data != size {
			return
		}
		process(buf, nodes)
	}
}

func process(data []byte, nodes redblacktree.Tree) error {
	operation := rune(data[0])
	if operation == 'I' {

	} else if operation == 'Q' {

	}
}

func insert(date int32, amount int32) {

}

func query(start int32, end int32) {

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

func read(conn net.Conn, buf []byte) (int, error) {
	n, err := conn.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		return n, err
	}
	return n, nil
}
