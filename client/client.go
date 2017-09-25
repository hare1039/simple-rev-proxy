package client

import (
	"fmt"
	"github.com/hare1039/simple-reverse-tunnel/def"
	"net"
	"os"
)

var Connections map[int]chan []byte

func backendHandler(backendTarget string, id int, toDemuxChan chan<- []byte) {
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
			fmt.Print("client backendHandler ")
			def.WriteConn(conn, bytes)
		case bytes := <-backendConn:
			fmt.Print("client backendHandler ")
			TCPs := def.TCPstream{
				Id:   id,
				Data: bytes,
			}
			toDemuxChan <- TCPs.Bytify()
		}
	}
}

func demuxConn(conn net.Conn, backendTarget string) {
	fromTunnel := make(chan []byte)
	fromBackendHandler := make(chan []byte, def.CHANNEL_BUF_AMOUNT)
	go def.ReadConnInJson(conn, fromTunnel)
	for {
		select {
		case rawBytes := <-fromTunnel:
			fmt.Print("clinet demuxConn ")
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
			fmt.Print("clinet demuxConn ")
			def.WriteConn(conn, bytes)
		}
	}
}

func Connect(target string, backend string) {
	conn, err := net.Dial("tcp", target)
	if err != nil {
		fmt.Println("Error connecting to:", err.Error())
		os.Exit(def.EXIT_ERROR_INTERNET)
	}
	defer conn.Close()

	Connections = make(map[int]chan []byte)
	demuxConn(conn, backend)
}
