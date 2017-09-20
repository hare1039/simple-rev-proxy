package def

import (
	"encoding/json"
	"fmt"
)

func (ts *TCPstream) Bytify() []byte {
	if bytes, err := json.Marshal(ts); err != nil {
		fmt.Println("Json bytify failed:", err.Error())
		return []byte("")
	} else {
		return bytes
	}
}

func ByteToTCPstream(b []byte) (TCPstream, bool) {
	var TCPs TCPstream
	if err := json.Unmarshal(b, &TCPs); err != nil {
		fmt.Println("ByteToTCPstream failed:", err.Error(), "string:", string(b))
		return TCPstream{}, false
	} else {
		return TCPs, true
	}

}
