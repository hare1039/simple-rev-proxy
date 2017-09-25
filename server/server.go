package server

import (
	"fmt"
	"github.com/hare1039/simple-reverse-tunnel/def"
	"net"
	"os"
)

var streams []chan def.TCPstream

func inboundHandler(conn net.Conn, inChan <-chan def.TCPstream) {
	clientSend := make(chan []byte)
	go def.ReadConnInJson(conn, clientSend)
	defer conn.Close()
	for {
		select {
		case CS := <-clientSend:
			fmt.Print("Server inboundHandler ")
			if TCPs, ok := def.ByteToTCPstream(CS); ok {
				streams[TCPs.Id] <- TCPs
			}
		case in := <-inChan:
			fmt.Print("Server inboundHandler ")
			def.WriteConn(conn, in.Bytify())
		}
	}
}

func outboundHandler(conn net.Conn, id int, forwarder chan<- def.TCPstream) int {
	outboundClientSend := make(chan []byte)
	go def.ReadConn(conn, outboundClientSend)
	defer conn.Close()
	for {
		select {
		case CS := <-streams[id]:
			fmt.Print("server outboundHandler ")
			def.WriteConn(conn, CS.Data)
		case in := <-outboundClientSend:
			fmt.Print("server outboundHandler ")
			//			fmt.Println("from outboundClientSend:", in)
			forwarder <- def.TCPstream{
				Id:   id,
				Data: in,
			}
		}
	}
	return id
}

func outboundServer(netInterface string, toInboundChan chan def.TCPstream) {
	fmt.Println("outboundServer")
	outln, err := net.Listen("tcp", netInterface)
	if err != nil {
		fmt.Println("Error listening to", err.Error())
		return
	} else {
		fmt.Println("Start reverse tunnel on: ", netInterface)
	}
	defer outln.Close()
	counter := 0
	for {
		if conn, err := outln.Accept(); err != nil {
			fmt.Println("Error on server accepting: ", err.Error())
		} else {
			streams = append(streams, make(chan def.TCPstream))
			fmt.Println("client: ", counter)
			go outboundHandler(conn, counter, toInboundChan)
			counter++
		}
	}
}

func Start(NetInterface string) {
	fmt.Println("starting server on " + NetInterface)
	ln, err := net.Listen("tcp", NetInterface)
	if err != nil {
		fmt.Println("Error on server listening: ", err.Error())
		os.Exit(2)
	}
	defer ln.Close()
	for {
		if conn, err := ln.Accept(); err != nil {
			fmt.Println("Error on server accepting: ", err.Error())
		} else {
			out2in := make(chan def.TCPstream, def.CHANNEL_BUF_AMOUNT)

			go inboundHandler(conn, out2in)
			go outboundServer("0.0.0.0:5421", out2in)
		}
	}
}
