package server

import (
	"fmt"
	"github.com/hare1039/simple-reverse-proxy/def"
	"net"
	"os"
)

var streams []chan def.TCPstream

func readConn(conn net.Conn, out chan<- []byte) {
	fmt.Println("readConn")
	buf := make([]byte, 32)
	for {
		reqLen, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading: ", err.Error())
			break
		} else if reqLen == 0 {
			break
		} else {
			fmt.Println("Sent buf")
			out <- buf
		}
	}
}

func writeConn(conn net.Conn, buf []byte) {
	fmt.Println("writeConn")
	conn.Write(buf)
}

func inboundHandler(conn net.Conn, inChan <-chan def.TCPstream) {
	fmt.Println("inboundHandler")
	clientSend := make(chan []byte)
	go readConn(conn, clientSend)
	defer conn.Close()
	for {
		select {
		case CS := <-clientSend:
			TCPs := def.ByteToTCPstream(CS)
			streams[TCPs.Id] <- TCPs
		case in := <-inChan:
			go writeConn(conn, in.Bytify())
		}
	}
}

func outboundHandler(conn net.Conn, id int, forwarder chan<- def.TCPstream) int {
	fmt.Println("outboundHandler")
	outboundClientSend := make(chan []byte)
	go readConn(conn, outboundClientSend)
	defer conn.Close()
	for {
		select {
		case CS := <-streams[id]:
			go writeConn(conn, CS.Bytify())
		case in := <-outboundClientSend:
			fmt.Println("data: ", in)
			S := def.TCPstream{
				Id:   id,
				Data: in,
			}
			forwarder <- S
		}
	}
	return id
}

func outboundServer(netInterface string, toInboundChan chan def.TCPstream) {
	fmt.Println("outboundServer")
	outln, err := net.Listen("tcp", netInterface)
	defer outln.Close()
	if err != nil {
		fmt.Println("Error listening to ", err.Error())
		return
	}
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
	defer ln.Close()
	if err != nil {
		fmt.Println("Error on server listening: ", err.Error())
		os.Exit(2)
	}
	for {
		if conn, err := ln.Accept(); err != nil {
			fmt.Println("Error on server accepting: ", err.Error())
		} else {
			out2in := make(chan def.TCPstream)

			go inboundHandler(conn, out2in)
			go outboundServer("0.0.0.0:5421", out2in)
		}
	}
}
