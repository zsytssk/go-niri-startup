package utils

import (
	"encoding/json"
	"fmt"
	"net"
)

func IsPortAvailable(port int) bool {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false
	}
	listener.Close()
	return true
}

type Msg = map[string]json.RawMessage

func GetData(msg []byte) (string, json.RawMessage) {
	m := make(Msg)
	json.Unmarshal(msg, &m)
	for k := range m {
		return k, m[k]
	}
	panic("")
}
