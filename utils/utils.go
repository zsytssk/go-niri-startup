package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
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

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return "", err
	}

	outBytes, _ := io.ReadAll(stdout)
	errBytes, _ := io.ReadAll(stderr)

	if err := cmd.Wait(); err != nil {
		return "", fmt.Errorf("%w: %s", err, string(errBytes))
	}

	return strings.TrimSpace(string(outBytes)), nil
}

func ReturnHttp(w http.ResponseWriter, msg string) {
	// 状态码
	w.WriteHeader(http.StatusOK)

	// Body
	w.Write([]byte(msg))
}

func GetCurDirFileName(fileSubFix string) (fileName string, err error) {
	ex, err := os.Executable()
	if err != nil {
		return
	}
	filepath.Base(ex)
	basename := filepath.Base(ex)
	fileName = fmt.Sprintf(`%s.%s`, basename, fileSubFix)
	return
}

func GetCurDirFilePath(fileName string) (filePath string, err error) {
	ex, err := os.Executable()
	if err != nil {
		return
	}
	filepath.Base(ex)
	exPath := filepath.Dir(ex)
	filePath = path.Join(exPath, fileName)
	return
}
