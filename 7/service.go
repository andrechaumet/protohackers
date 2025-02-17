package main

import (
	"sync"
	"time"
)

var Types = map[string]func(Request){
	"ack":     ack,
	"data":    data,
	"close":   cls,
	"connect": conn,
}

type Request struct {
	tp   string
	id   int
	pos  int
	msg  string
	time time.Time
}

type sequence struct {
	reqs []*Request
	lock *sync.Mutex
}

var sequences = make(map[int]*sequence, 10)

func conn(req Request) {
	if val := sequences[req.id]; val != nil {
		return
	} else {
		sequences[req.id] = &sequence{make([]*Request, 10), nil}
	}
}

func data(req Request) {
	if val := sequences[req.id]; val == nil {
		return
	} else {
		val.reqs = append(val.reqs, &req)
	}
}

func ack(req Request) {
	/*
		If the SESSION is not open: send /close/SESSION/ and stop.
		If the LENGTH value is not larger than the largest LENGTH value in any ack message you've received on this session so far: do nothing and stop (assume it's a duplicate ack that got delayed).
		If the LENGTH value is larger than the total amount of payload you've sent: the peer is misbehaving, close the session.
		If the LENGTH value is smaller than the total amount of payload you've sent: retransmit all payload data after the first LENGTH bytes.
		If the LENGTH value is equal to the total amount of payload you've sent: don't send any reply.
	*/
}

func cls(req Request) {
	/*
		server will never need to initiate the closing of any sessions.
		/close/SESSION/ message, send a matching close message back.
	*/
}
