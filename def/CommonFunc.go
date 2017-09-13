package def

import (
	"fmt"
	"net"
)

func ReadConn(conn net.Conn, out chan<- []byte) {
	for {
		buf := make([]byte, BUF_SIZE)
		reqLen, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading: ", err.Error())
			break
		} else if reqLen != 0 {
			out <- buf
		}
	}
}

func WriteConn(conn net.Conn, buf []byte) {
	conn.Write(buf)
}
