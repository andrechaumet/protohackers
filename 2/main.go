package main

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"sync"
)

/*
	protohackers.com/problem/2

	byte:  |  0   |  1     2     3     4  |  5     6     7     8  |
	type:  | char |         int32         |         int32         |
*/

type insertion struct {
	time   int32
	amount int32
}

const size = 9

var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, size)
	},
}

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
	inserts := make([]insertion, 0)
	for {
		buf := [9]byte{}
		data, err := read(conn, buf)
		if err != nil || data != size {
			return
		}
		if sum, query := process(buf, inserts); query {
			sumBytes := make([]byte, 4)
			binary.BigEndian.PutUint32(sumBytes, uint32(sum))
			conn.Write(sumBytes)
		}
	}
}

func process(data [9]byte, inserts []insertion) (int32, bool) {
	operation := rune(data[0])
	if operation == 'I' {
		insert(data, &inserts)
	} else if operation == 'Q' {
		return query(data, &inserts), true
	}
	return 0, false
}

func insert(data [9]byte, insertions *[]insertion) {
	time := convert(data[1:5])
	price := convert(data[6:9])
	*insertions = append(*insertions, insertion{time, price})
	log.Println("inserted %d at %d\n", price, time)
}

func query(data [9]byte, insertions *[]insertion) int32 {
	start := convert(data[1:5])
	end := convert(data[6:9])
	var sum int32
	for _, insertion := range *insertions {
		if insertion.time > start && insertion.time < end {
			sum += insertion.amount
		}
	}
	log.Println("sum from %d to %d\n", start, end)
	return sum
}

func convert(data []byte) int32 {
	return int32(data[0])<<24 | int32(data[1])<<16 | int32(data[2])<<8 | int32(data[3])
}

func read(conn net.Conn, buf [9]byte) (int, error) {
	n, err := conn.Read(buf) //aaaaa ayuda
	if err != nil && !errors.Is(err, io.EOF) {
		return n, err
	}
	return n, nil
}
