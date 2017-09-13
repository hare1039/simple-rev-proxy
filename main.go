package main

import (
	"fmt"
	"github.com/hare1039/simple-reverse-tunnel/client"
	"github.com/hare1039/simple-reverse-tunnel/server"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify `simple-rev-tunnel server` or `simple-rev-tunnel client`")
		os.Exit(1)
	}
	if os.Args[1] == "server" || os.Args[1] == "s" {
		server.Start(":1256")
	} else if os.Args[1] == "client" || os.Args[1] == "c" {
		client.Connect("localhost:1256")
	}
}
