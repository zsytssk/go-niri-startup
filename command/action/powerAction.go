package action

import (
	"niri-startup/state"
	"niri-startup/utils"
	"time"
)

func PowerAction() error {
	result, err := utils.RunCMD(`printf "󰌾 Lock\n󰍃 Logout\n󰙧 Shutdown\n󰑐 Reboot\n󰚰 Update" | fuzzel -d -p "请选择: "`, false)
	if err != nil {
		return err
	}
	if result == "󰌾 Lock" {
		utils.RunCMD("swaylock --daemonize", false)
		time.Sleep(1 * time.Second)
		utils.NiriSendAction(utils.Action{
			PowerOffMonitors: &utils.Empty{},
		})
		return nil
	}
	if result == "󰍃 Logout" {
		utils.RunCMD("niri msg action quit --skip-confirmation", false)
		return nil
	}
	if result == "󰑐 Reboot" {
		utils.RunCMD(`zenity --question --text="确定要重启吗？"`, false)
		utils.RunCMD("reboot", false)
		return nil
	}
	if result == "󰙧 Shutdown" {
		utils.RunCMD(`zenity --question --text="确定要关机吗？"`, false)
		utils.RunCMD("shutdown -h now", false)
		return nil
	}
	if result == "󰚰 Update" {
		utils.RunCMD(
			`ghostty --title="Update System" --class="update.ghostty" -e sh -c "neofetch && sudo apt update && sudo apt upgrade; exec bash"`,
			true,
		)
		instance := state.GetStateInstance()
		waitWindowOpen := state.UseWaitWindowOpen(instance)
		item, err := waitWindowOpen(func(w *state.Window) bool {
			return w.AppId == "update.ghostty"
		})
		if err != nil {
			return err
		}
		currentWorkspaceId := instance.CurrentWorkspaceId
		utils.NiriSendActionArr([]utils.Action{
			{
				MoveWindowToWorkspace: &utils.MoveWindowToWorkspace{
					WindowId:  item.ID,
					Focus:     true,
					Reference: utils.WindowReference{Id: currentWorkspaceId},
				},
			},
			{
				SetWindowHeight: &utils.SetWindowSize{Id: item.ID,
					Change: utils.SetWindowSizeChange{SetFixed: 900},
				},
			},
			{
				SetWindowWidth: &utils.SetWindowSize{Id: item.ID,
					Change: utils.SetWindowSizeChange{SetFixed: 900},
				},
			},
			{
				MoveWindowToFloating: &utils.WindowWithId{Id: item.ID},
			},
			{Sleep: 80},
			{FocusWindow: &utils.WindowWithId{Id: item.ID}},
			{
				CenterWindow: &utils.WindowWithId{Id: item.ID},
			},
		})
	}
	return nil
}
