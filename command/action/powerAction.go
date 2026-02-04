package action

import (
	"fmt"
	"niri-startup/action"
	"niri-startup/state"
	"niri-startup/utils"
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

func PowerAction() error {
	result, err := utils.RunCMD(`printf "󰌾 Lock\n󰍃 Logout\n󰙧 Shutdown\n󰑐 Reboot\n󰚰 Update" | fuzzel -d -p "请选择: "`, false)
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
		utils.RunCMD("niri msg action quit --skip-confirmation", false)
		return nil
	}
	if result == "󰑐 Reboot" {
		if runConfirm(`确定要重启吗？`) {
			utils.RunCMD("reboot", false)
		}
		return nil
	}
	if result == "󰙧 Shutdown" {
		if runConfirm(`确定要关机吗？`) {
			utils.RunCMD("shutdown -h now", false)
		}
		return nil
	}
	if result == "󰚰 Update" {
		utils.RunCMD(
			`ghostty --title="Update System" --class="update.ghostty" -e sh -c "neofetch && sudo apt update && sudo apt upgrade; exec bash"`,
			true,
		)
		instance := state.GetStateInstance()
		waitWindowOpen := state.UseWaitWindowOpen(instance)
		view := instance.GetSnapshot()
		item, err := waitWindowOpen(func(w *state.Window) bool {
			return w.AppId == "update.ghostty"
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
	}
	return nil
}
