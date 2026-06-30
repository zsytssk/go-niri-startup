package cmdAction

import (
	"fmt"
	"os/exec"
	"strings"
)

var RelaunchAction = CmdActionItem {
	CmdList: []string {
		"Waybar Relaunch",
		"Transws Relaunch",
		"Niri-startup Relaunch",
	},
	Fn: CmdAction,
}

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

func CmdAction(cmd string) error {
	var err error
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
