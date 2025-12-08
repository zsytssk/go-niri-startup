package main

import (
	"bytes"
	"fmt"
	"net/http"
	"niri-startup/command"
	"niri-startup/command/action"
	"niri-startup/command/spad"
	"niri-startup/config"
	"niri-startup/state"
	"niri-startup/utils"
	"os"
)

func getCmd() []string {
	if len(os.Args) == 1 {
		return nil
	}
	return os.Args[1:]
}

const PORT = 6321

func main() {
	// 处理命令行参数
	args := getCmd()
	if len(args) > 0 {
		name, data := args[0], args[1]
		http.Post(fmt.Sprintf("http://127.0.0.1:%d/%s", PORT, name), "application/json", bytes.NewBuffer([]byte(data)))
		return
	}
	// 处理本地服务器
	if !utils.IsPortAvailable(PORT) {
		panic(fmt.Sprintf("端口 %d 已被占用", PORT))
	}
	_, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	state.GetStateInstance()
	utils.GetSocketInstance()
	utils.RunCMD("notify-send 启动 niri-ts-startup!", false)
	http.HandleFunc("/spad", spad.Spad)
	http.HandleFunc("/action", action.Action)
	http.HandleFunc("/runApp", command.RunApp)
	http.HandleFunc("/getState", command.GetState)
	err = http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
	if err != nil {
		panic(err)
	}
}
