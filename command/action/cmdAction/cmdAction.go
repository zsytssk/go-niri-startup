package cmdAction

import (
	"fmt"
	"niri-startup/utils"
	"strings"
)

type CmdActionItem struct {
	CmdList  []string
	Fn     func(string)error
}

func Run(prefix string) error {
	actionList := []CmdActionItem {
		RelaunchAction,
		HdmiOutoutPositionAction,
		SystemAction,
	}
	cmdList := make([]string, 0)
	for _,item := range(actionList) {
		cmdList = append(cmdList, item.CmdList...)
	}
	var cmd string
	if len(prefix) > 0 {
		cmd = fmt.Sprintf(`printf "%s" | fuzzel -d --search "%s" -p "快捷命令: "`, strings.Join(cmdList, "\n"), prefix)
	} else {
		cmd = fmt.Sprintf(`printf "%s" | fuzzel -d -p "快捷命令: "`, strings.Join(cmdList, "\n"))
	}
	result, err := utils.RunCMD(cmd, false)
	// fmt.Println(cmd, err)
	if err != nil {
		return err
	}
	for _,item := range(actionList) {
		if utils.Contains(item.CmdList, result) {
			err = item.Fn(result)
			if err != nil {
				return err
			}
			break
		}
	}

	return nil
}
