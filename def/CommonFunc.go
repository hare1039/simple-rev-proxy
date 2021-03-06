package def

import (
	"bytes"
	"fmt"
	"net"
	"time"
)

func ReadConn(conn net.Conn, out chan<- []byte) {
	fmt.Println("ReadConn")
	for {
		buf := make([]byte, BUF_SIZE)
		if n, err := conn.Read(buf); err != nil {
			fmt.Println("Error reading: ", err.Error())
			break
		} else if n > 0 {
			fmt.Println(time.Now(), "Readed buffer:", string(buf))
			out <- buf
		}
	}
}

func ReadConnInJson(conn net.Conn, out chan<- []byte) {
	buf := make([]byte, BUF_SIZE)
	fmt.Println("ReadConn")
	var jsonBuf bytes.Buffer
	stack := 0
	for {
		if n, err := conn.Read(buf); err != nil {
			fmt.Println("Error reading: ", err.Error())
			break
		} else if n > 0 {
			limit := n
			for i, x := range buf[:n] {
				if x == []byte("{")[0] {
					stack++
				} else if x == []byte("}")[0] {
					stack--
					limit = i
					break
				}
			}
			if stack == 0 {
				jsonBuf.Write(buf[:limit+1])
				var sender []byte
				sender = append(sender, jsonBuf.Bytes()...)
				out <- sender
				fmt.Println(time.Now(), "send to channel:", string(sender))
				jsonBuf.Reset()
				for i, x := range buf[limit+1:] {
					if x == []byte("{")[0] {
						fmt.Println("append left over:", string(buf[limit+i+1:]))
						jsonBuf.Write(buf[limit+i+1:])
						stack++
						break
					}
				}
			} else {
				jsonBuf.Write(buf[:n])
				fmt.Println("appended buffer:", string(jsonBuf.Bytes()))
			}
		}
	}
}

func WriteConn(conn net.Conn, buf []byte) {
	fmt.Println(time.Now(), "write connection:", string(buf))
	conn.Write(buf)
}
