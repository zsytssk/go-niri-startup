package action

import "niri-startup/utils"

func ToggleInput() {
	cur, _ := utils.RunCMD("fcitx5-remote", false)
	if cur == "2" {
		utils.RunCMD("fcitx5-remote -c", false)
	} else {
		utils.RunCMD("fcitx5-remote -o", false)
	}
}
