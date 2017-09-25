package main

import (
	"fmt"
	"github.com/hare1039/simple-reverse-tunnel/client"
	"github.com/hare1039/simple-reverse-tunnel/server"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Please specify `simple-reverse-tunnel server [:port]` or `simple-reverse-tunnel client [ip:port] [ip:port]`")
		os.Exit(1)
	}
	if os.Args[1] == "server" || os.Args[1] == "s" {
		server.Start(os.Args[2])
	} else if os.Args[1] == "client" || os.Args[1] == "c" {
		client.Connect(os.Args[2], os.Args[3])
	}
}
