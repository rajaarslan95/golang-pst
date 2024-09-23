package main

import (
	"flag"
	"user-manager/rest"
)

var addr = flag.String("a", "127.0.0.1", "Address to listen")
var port = flag.Int("p", 9001, "Port to listen")

func main() {
	flag.Parse()

	server := rest.NewServer(*addr, *port)

	server.Start()
}
