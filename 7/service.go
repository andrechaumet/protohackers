package main

var Types = map[string]func(Request){
	"connect": conn,
	"data":    data,
	"ack":     ack,
	"close":   cls,
}

type Request struct {
	tp  string
	id  int
	pos int
	msg string
}

var requests = make(map[string]Request, 10)

func conn(req Request) {
}

func data(req Request) {

}

func ack(req Request) {
}

func cls(req Request) {

}
