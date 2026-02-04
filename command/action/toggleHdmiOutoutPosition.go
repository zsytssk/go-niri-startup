package action

import (
	"niri-startup/utils"
)

func toggleHdmiOutoutPosition() {
	// utils.NiriSendOutputAction(utils.OutputAction{
	// 	Output: "HDMI-A-1",
	// 	Action: &utils.OutputActionCon{
	// 		Off: &null,
	// 	},
	// })
	outputsStr, err := utils.RunCMD("niri msg --json outputs", false)
	if err != nil {
		return
	}
	posX, err := utils.GetNestedValue(outputsStr, "HDMI-A-1.logical.x")
	if err != nil {
		return
	}
	if posX.(float64) > 0 {
		utils.RunCMD("niri msg output HDMI-A-1 position set 0 0", false)
	} else {
		utils.RunCMD("niri msg output HDMI-A-1 position set 1920 0", false)
	}
	utils.RunCMD("notify-send 'Toggle HDMI-A-1 Position'", false)
	// utils.NiriSendOutputAction(utils.OutputAction{
	// 	Output: "HDMI-A-1",
	// 	Action: &utils.OutputActionCon{
	// 		Position: &utils.OutputActionPosition{
	// 			// Automatic: &utils.Empty{},
	// 			Specific: &utils.OutputActionPositionSpecific{
	// 				X: 1980,
	// 				Y: 0,
	// 			},
	// 		},
	// 	},
	// })
}
