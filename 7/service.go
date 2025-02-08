package main

var Types = map[string]func([]string){
	"connect": conn,
	"data":    data,
	"ack":     ack,
	"close":   cls,
}

var connections = map[string][]string{}

func conn(args []string) {
}

func data(args []string) {
}

func ack(args []string) {
}

func cls(args []string) {
}
