package def

import (
	"fmt"
	"net"
)

func ReadConn(conn net.Conn, out chan<- []byte) {
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
			fmt.Println("Sent buf:", buf)
			out <- buf
		}
	}
}

func WriteConn(conn net.Conn, buf []byte) {
	fmt.Println("writeConn")
	conn.Write(buf)
}
