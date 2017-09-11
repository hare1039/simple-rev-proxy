package main

import (
	"fmt"
	"github.com/hare1039/simple-reverse-proxy/server"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify `simple-rev-proxy server` or `simple-rev-proxy client`")
		os.Exit(1)
	}
	if os.Args[1] == "server" || os.Args[1] == "s" {
		server.Start(":1256")
	}
}
