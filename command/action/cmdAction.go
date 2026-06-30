package action

import (
	"fmt"
	"niri-startup/utils"
	"os/exec"
	"strings"
)

func runNiriSh(sh string) error {
	// niri msg action spawn-sh -- "你要执行的命令"

	cmd := exec.Command(
		"niri",
		"msg",
		"action",
		"spawn-sh",
		"--",
		sh,
	)

	err := cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

func CmdAction() error {

	cmd, err := utils.RunCMD(`printf "Waybar Relaunch\nTransws Relaunch\nNiri-startup Relaunch\n" | fuzzel -d -p "运行快捷命令: "`, false)
	if err != nil {
		return err
	}
	if cmd == "Niri-startup Relaunch" {
		err = runNiriSh("pkill niri-startup >/dev/null 2>&1 || true; exec niri-startup >> ~/Documents/zsy/github/go-niri-startup/niri-bun.log 2>&1")
		if err != nil {
			return err
		}
		return nil
	}

	cmd = strings.ToLower(cmd)
	cmd = fmt.Sprintf("(pkill %s 2>/dev/null || true) && %s &", cmd, cmd)
	fmt.Println(cmd)
	err = runNiriSh(cmd)
	if err != nil {
		return err
	}
	return nil
}
