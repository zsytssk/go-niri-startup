package main

import (
	"bytes"
	"fmt"
	"net/http"
	"niri-startup/command"
	"niri-startup/config"
	"niri-startup/utils"
	"os"
)

func getCmd() []string {
	if len(os.Args) == 1 {
		return nil
	}
	return os.Args[1:]
}

const port = 6322

func main() {
	// 处理命令行参数
	args := getCmd()
	if len(args) > 0 {
		name, data := args[0], args[1]
		http.Post(fmt.Sprintf("http://127.0.0.1:6322/%s", name), "application/json", bytes.NewBuffer([]byte(data)))
		return
	}
	fmt.Println(args)

	// 处理本地服务器
	if !utils.IsPortAvailable(port) {
		panic(fmt.Sprintf("端口 %d 已被占用", port))
	}
	_, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/spad", command.Spad)
	http.HandleFunc("/action", command.Action)
	http.HandleFunc("/runApp", command.RunApp)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic(err)
	}

}
