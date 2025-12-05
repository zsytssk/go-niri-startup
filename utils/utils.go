package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"strings"
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
func GetMsgType(msg []byte) string {
	m := make(Msg)
	json.Unmarshal(msg, &m)
	for k := range m {
		return k
	}
	return ""
}

func RunCMD(input string, nohup bool) (string, error) {
	var cmd *exec.Cmd
	if nohup {
		cmd = exec.Command("bash", "-c",
			fmt.Sprintf("nohup sh -c '%s' > /dev/null 2>&1 &", input),
		)
	} else {
		cmd = exec.Command("bash", "-c", input)
	}

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf // 如果你也想捕获错误输出

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}

func ReturnHttp(w http.ResponseWriter, msg string) {
	// 状态码
	w.WriteHeader(http.StatusOK)

	// Body
	w.Write([]byte(msg))
}
