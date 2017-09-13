package client

import (
	"fmt"
	"github.com/hare1039/simple-reverse-tunnel/def"
	"net"
)

var Connections map[int]chan []byte

func backendHandler(backendTarget string, id int, toDemuxChan chan<- []byte) {
	fmt.Println("backendHandler")
	conn, err := net.Dial("tcp", backendTarget)
	if err != nil {
		fmt.Println("Error connecting to backend:", err.Error())
	}
	defer conn.Close()
	fmt.Println("Connect to: ", backendTarget)
	backendConn := make(chan []byte)
	go def.ReadConn(conn, backendConn)
	for {
		select {
		case bytes := <-Connections[id]:
			def.WriteConn(conn, bytes)
		case bytes := <-backendConn:
			TCPs := def.TCPstream{
				Id:   id,
				Data: bytes,
			}
			toDemuxChan <- TCPs.Bytify()
		}
	}
}

func demuxConn(conn net.Conn, backendTarget string) {
	fmt.Println("demuxConn")
	var fromTunnel chan []byte
	var fromBackendHandler chan []byte
	go def.ReadConn(conn, fromTunnel)
	for {
		select {
		case rawBytes := <-fromTunnel:
			fmt.Println(rawBytes)
			TCPs := def.ByteToTCPstream(rawBytes)
			if ch, ok := Connections[TCPs.Id]; ok {
				ch <- TCPs.Data
			} else {
				ch := make(chan []byte)
				Connections[TCPs.Id] = ch
				go backendHandler(backendTarget, TCPs.Id, fromBackendHandler)
				ch <- TCPs.Data
			}
		case bytes := <-fromBackendHandler:
			def.WriteConn(conn, bytes)
		}
	}
}

func Connect(target string) {
	conn, err := net.Dial("tcp", target)
	if err != nil {
		fmt.Println("Error connecting to:", err.Error())
	}
	defer conn.Close()

	demuxConn(conn, "localhost:9988")
}
