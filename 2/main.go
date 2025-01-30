package main

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
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

const protoSize = 9

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
		log.Println("Received a new connection")
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	log.Println("Starting new process")
	inserts := make([]insertion, 0)
	for {
		var buf = make([]byte, 9)
		size, err := read(conn, buf)
		if err != nil || size != protoSize {
			return
		}
		if sum, query := process(buf, &inserts); query {
			sumBytes := make([]byte, 4)
			binary.BigEndian.PutUint32(sumBytes, uint32(sum))
			conn.Write(sumBytes)
		}
	}
}

func process(data []byte, inserts *[]insertion) (int32, bool) {
	operation := rune(data[0])
	if operation == 'I' {
		insert(data, inserts)
	} else if operation == 'Q' {
		return query(data, inserts), true
	} else {
		log.Println("invalid operation", operation)
	}
	return 0, false
}

func insert(data []byte, insertions *[]insertion) {
	time := convert(data[1:5])
	price := convert2(data[6:9])
	*insertions = append(*insertions, insertion{time, price})
	log.Println("time: ", time, " ", " price: ", price)
}

func query(data []byte, insertions *[]insertion) int32 {
	start := convert(data[1:5])
	end := convert2(data[6:9])
	var sum int32
	for _, insertion := range *insertions {
		if insertion.time > start && insertion.time < end {
			sum += insertion.amount
		}
	}
	log.Println("sum: ", sum)
	return sum
}

func convert(data []byte) int32 {
	if len(data) < 4 {
		data = append(data, make([]byte, 4-len(data))...)
	}
	con := int32(data[0])<<24 | int32(data[1])<<16 | int32(data[2])<<8 | int32(data[3])
	return con
}

// ok so im still figuring this one out
func convert2(data []byte) int32 {
	if len(data) < 4 {
		data = append(data, make([]byte, 4-len(data))...)
	}
	var result int32
	length := len(data)

	if length > 0 {
		result |= int32(data[0]) << 16
	}
	if length > 1 {
		result |= int32(data[1]) << 8
	}
	if length > 2 {
		result |= int32(data[2])
	}
	return result
}

/* this fella is in jail until new notice
func convert(data []byte) int32 {
	var result int32
	for i, b := range data {
		result |= int32(b) << (24 - 8*i)
	}
	return result
}*/

func read(conn net.Conn, buf []byte) (int, error) {
	n, err := conn.Read(buf[:])
	if err != nil && !errors.Is(err, io.EOF) {
		return n, err
	}
	return n, nil
}
