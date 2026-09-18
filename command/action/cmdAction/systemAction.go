package cmdAction

import (
	"fmt"
	"niri-startup/action"
	"niri-startup/state"
	"niri-startup/utils"
	"os/exec"
	"strings"
)

var SystemAction = CmdActionItem {
	CmdList: []string {
		"Update System",
		"KillWine System",
	},
	Fn: SystemActionFn,
}


func SystemActionFn(cmd string) error {
	var err error
	cmd = strings.Replace(cmd, " System", "", 1)
	if cmd == "Update" {
		err = runGhosttyCmd("update", "Update System", "neofetch && sudo apt update && sudo apt upgrade; exec bash")
		if err != nil {
			return err
		}
	}
	if cmd == "KillWine" {
		err = runPowerOption("KillWine")
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}


func runConfirm(tip string) bool {
	cmd := fmt.Sprintf(`if zenity --question --text="%s" --title="问题"; then echo "Y"; else echo "N"; fi`, tip)
	r, err := utils.RunCMD(cmd, false)
	if err != nil || r == "N" {
		return false
	}
	return true
}

func runGhosttyCmd(name string, title string, script string) error {
	// cmd = fmt.Sprintf(`ghostty --title="%s" --class="runCmd.ghostty" -e sh -c "%s"`, title, cmd)
	cmd := exec.Command(
		"ghostty",
		"--title="+title,
		"--class="+name+".ghostty",
		"-e",
		"sh",
		"-c",
		script,
	)

	err := cmd.Start()
	if err != nil {
		return err
	}
	instance := state.GetStateInstance()
	waitWindowOpen := state.UseWaitWindowOpen(instance)
	waitWindowClose := state.UseWaitWindowClose(instance)
	view := instance.GetSnapshot()
	item, err := waitWindowOpen(func(w *state.Window) bool {
		return w.AppId == name+".ghostty"
	})
	if err != nil {
		return err
	}
	currentWorkspaceId := view.CurrentWorkspaceId
	utils.NiriSendActionArr([]action.Action{
		{
			MoveWindowToWorkspace: &action.MoveWindowToWorkspace{
				WindowId:  item.ID,
				Focus:     true,
				Reference: action.WindowReference{Id: currentWorkspaceId},
			},
		},
		{
			SetWindowHeight: &action.SetWindowSize{Id: item.ID,
				Change: action.SetWindowSizeChange{SetFixed: 900},
			},
		},
		{
			SetWindowWidth: &action.SetWindowSize{Id: item.ID,
				Change: action.SetWindowSizeChange{SetFixed: 900},
			},
		},
		{
			MoveWindowToFloating: &action.WindowWithId{Id: item.ID},
		},
		{Sleep: 80},
		{FocusWindow: &action.WindowWithId{Id: item.ID}},
		{
			CenterWindow: &action.WindowWithId{Id: item.ID},
		},
	})
	err = waitWindowClose(item)
	if err != nil {
		return err
	}
	return nil
}

func runPowerOption(name string) error {
	script := fmt.Sprintf(`
		clear

		echo "===== %s ====="

		if pgrep -f "wine|wineserver|\.exe" > /dev/null 2>&1; then
			echo "[关闭] Wine..."
			/usr/bin/wineserver -k
			echo "[✓] Wine 已关闭"
		fi

		if pgrep -x postgres >/dev/null 2>&1; then
			echo "[关闭] Postgres..."
			sudo pkill -x postgres
			echo "[✓] Postgres 已关闭"
		fi

		echo
		exec bash
	`, name)

	err := runGhosttyCmd("clean", "Clean System", script)
	if err != nil {
		return err
	}
	return nil
}
