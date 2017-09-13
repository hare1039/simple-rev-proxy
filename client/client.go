package client

import (
	"fmt"
	"github.com/hare1039/simple-reverse-tunnel/def"
	"net"
	"os"
)

var Connections map[int]chan []byte

func backendHandler(backendTarget string, id int, toDemuxChan chan<- []byte) {
	fmt.Println("backendHandler")
	conn, err := net.Dial("tcp", backendTarget)
	if err != nil {
		fmt.Println("Error connecting to backend:", err.Error())
		return
	}
	defer conn.Close()
	fmt.Println("Connected to:", backendTarget)
	backendConn := make(chan []byte)
	go def.ReadConn(conn, backendConn)
	for {
		select {
		case bytes := <-Connections[id]:
			go def.WriteConn(conn, bytes)
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
	fromTunnel := make(chan []byte)
	fromBackendHandler := make(chan []byte)
	go def.ReadConn(conn, fromTunnel)
	for {
		select {
		case rawBytes := <-fromTunnel:
			if TCPs, Convertok := def.ByteToTCPstream(rawBytes); Convertok {
				if ch, ok := Connections[TCPs.Id]; ok {
					ch <- TCPs.Data
				} else {
					ch := make(chan []byte)
					Connections[TCPs.Id] = ch
					go backendHandler(backendTarget, TCPs.Id, fromBackendHandler)
					ch <- TCPs.Data
				}
			}
		case bytes := <-fromBackendHandler:
			go def.WriteConn(conn, bytes)
		}
	}
}

func Connect(target string) {
	conn, err := net.Dial("tcp", target)
	if err != nil {
		fmt.Println("Error connecting to:", err.Error())
		os.Exit(def.EXIT_ERROR_INTERNET)
	}
	defer conn.Close()

	Connections = make(map[int]chan []byte)
	demuxConn(conn, "localhost:9988")
}
