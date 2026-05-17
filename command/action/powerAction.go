package action

import (
	"fmt"
	"niri-startup/action"
	"niri-startup/state"
	"niri-startup/utils"
	"os/exec"
	"time"
)

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
		panic(err)
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

func runPowerOption(name string, cmd string) error {
	script := fmt.Sprintf(`
		clear

		echo "===== %s ====="

		if pgrep -x wineserver >/dev/null 2>&1 || pgrep -x wine >/dev/null 2>&1; then
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
		echo "是否%s? (Y/n)"
		read -p "> " choice

		choice="${choice:-Y}"

		case "$choice" in
			Y|y)
				echo "正在%s..."
				%s
				;;
			*)
				echo "已取消%s"
				;;
		esac
		exec bash
	`, name, name, name, cmd, name)

	err := runGhosttyCmd("clean", "Clean System", script)
	if err != nil {
		return err
	}
	return nil
}

func PowerAction() error {

	result, err := utils.RunCMD(`printf "󰌾 Lock\n󰍃 Logout\n󰙧 Shutdown\n󰑐 Reboot\n󰚰 Update\n☕ Test" | fuzzel -d -p "请选择: "`, false)
	if err != nil {
		return err
	}
	if result == "󰌾 Lock" {
		utils.RunCMD("swaylock --daemonize", true)
		time.Sleep(1 * time.Second)
		utils.NiriSendAction(action.Action{
			PowerOffMonitors: &action.Empty{},
		})
		return nil
	}
	if result == "󰍃 Logout" {
		err = runPowerOption("退出登陆", "niri msg action quit --skip-confirmation")
		if err != nil {
			return err
		}
		return nil
	}
	if result == "󰑐 Reboot" {
		err = runPowerOption("重启", "reboot")
		if err != nil {
			return err
		}
		return nil
	}
	if result == "󰙧 Shutdown" {
		err = runPowerOption("关机", "shutdown -h now")
		if err != nil {
			return err
		}
		return nil
	}
	if result == "☕ Test" {
		// err = runGhosttyCmd("Update System", "/usr/bin/wineserver -k ; sudo pkill -9 postgres; exec bash")
		// err = runPowerOption("关机", "shutdown -h now")
		err = runPowerOption("关机", "echo \"假装关机...\"")
		if err != nil {
			return err
		}
		return nil
	}
	if result == "󰚰 Update" {
		err = runGhosttyCmd("update", "Update System", "neofetch && sudo apt update && sudo apt upgrade; exec bash")
		if err != nil {
			return err
		}
	}
	return nil
}
