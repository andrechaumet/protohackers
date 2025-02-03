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
	ln := setup(":8082")
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
		log.Println("Invalid operation", operation)
	}
	return 0, false
}

func insert(data []byte, insertions *[]insertion) {
	time, price := convert(data[1:5]), convert(data[6:9])
	*insertions = append(*insertions, insertion{time, price})
	log.Printf("Inserted time: %v, price: %v", time, price)
}

func query(data []byte, insertions *[]insertion) int32 {
	start, end := convert(data[1:5]), convert(data[6:9])
	var sum int32
	for _, insertion := range *insertions {
		if insertion.time > start && insertion.time < end {
			sum += insertion.amount
		}
	}
	return sum
}

func convert(data []byte) int32 {
	for {
		if len(data) == 4 {
			break
		} else {
			data = append(data, 0)
			copy(data[1:], data[:len(data)-1])
			data[0] = 0
		}
	}
	return int32(binary.BigEndian.Uint32(data))
}

func read(conn net.Conn, buf []byte) (int, error) {
	n, err := conn.Read(buf[:])
	if err != nil && !errors.Is(err, io.EOF) {
		return n, err
	}
	return n, nil
}
