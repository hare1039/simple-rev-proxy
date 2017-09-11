package def

import (
	"fmt"
	"strconv"
)

func (ts *TCPstream) Bytify() []byte {
	return append([]byte(strconv.Itoa(ts.Id)+"|"), ts.Data...)
}

func ByteToTCPstream(b []byte) TCPstream {
	id := ""
	for i, _ := range b {
		if string(b[i]) != "|" {
			id = id + string(b[i])
		} else {
			if D, err := strconv.Atoi(id); err != nil {
				fmt.Println("Error on converting id: ", err.Error(), id)
			} else {
				return TCPstream{
					Id:   D,
					Data: b[i+1:],
				}
			}
		}
	}
	return TCPstream{
		Id:   -1,
		Data: []byte(""),
	}
}
