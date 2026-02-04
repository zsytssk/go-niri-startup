package action

import (
	"niri-startup/action"
	"niri-startup/utils"
)

func toggleHdmiOutoutPosition() {
	outputsStr, err := utils.NiriSend(map[string]interface{}{
		"Outputs": nil,
	})
	if err != nil {
		return
	}

	posX, err := utils.GetNestedValueFromStr(string(outputsStr), "Ok.Outputs.HDMI-A-1.logical.x")
	if err != nil {
		return
	}
	if posX.(float64) > 0 {
		utils.NiriSend(
			map[string]interface{}{
				"Output": action.OutputAction{
					Output: "HDMI-A-1",
					Action: &action.OutputActionCon{
						Position: &action.OutputActionPosition{
							Position: &action.OutputActionPositionDetail{
								Specific: &action.OutputActionPositionSpecific{X: 0, Y: 0},
							},
						}},
				},
			},
		)
	} else {
		utils.NiriSend(
			map[string]interface{}{
				"Output": action.OutputAction{
					Output: "HDMI-A-1",
					Action: &action.OutputActionCon{
						Position: &action.OutputActionPosition{
							Position: &action.OutputActionPositionDetail{
								Specific: &action.OutputActionPositionSpecific{X: 1920, Y: 0},
							},
						}},
				},
			},
		)
	}
}

// func toggleHdmiOutoutPosition() {
// utils.NiriSend(
// 	map[string]interface{}{
// 		"Output": map[string]interface{}{
// 			"output": "HDMI-A-1",
// 			"action": map[string]interface{}{"On": nil},
// 		},
// 	},
// )

// utils.NiriSend(
// 	map[string]interface{}{
// 		"Output": map[string]interface{}{
// 			"output": "HDMI-A-1",
// 			"action": map[string]interface{}{
// 				"Position": map[string]interface{}{
// 					"position": map[string]interface{}{
// 						"Specific": map[string]interface{}{"x": 1920, "y": 0},
// 					},
// 				}},
// 		},
// 	},
// )
// utils.NiriSend(
// 	map[string]interface{}{
// 		"Output": utils.OutputAction{
// 			Output: "HDMI-A-1",
// 			Action: &utils.OutputActionCon{
// 				Position: &utils.OutputActionPosition{
// 					Position: &utils.OutputActionPositionDetail{
// 						Specific: &utils.OutputActionPositionSpecific{X: 1920, Y: 0},
// 					},
// 				}},
// 		},
// 	},
// )
// utils.NiriSend(map[string]interface{}{
// 	"Version": nil,
// })
// outputsStr, err := utils.NiriSend(map[string]interface{}{
// 	"Outputs": nil,
// })

// if err != nil {
// 	return
// }
// utils.NiriSend(map[string]interface{}{
// 	"Action": map[string]interface{}{
// 		"Version": nil,
// 	}})
// utils.NiriSend(map[string]interface{}{
// 	"Action": map[string]interface{}{
// 		"ToggleWindowFloating": map[string]interface{}{"id": 9},
// 	}})

// outputsStr, err := utils.RunCMD("niri msg --json outputs", false)
// if err != nil {
// 	return
// }
// posX, err := utils.GetNestedValue(string(outputsStr), "Ok.Outputs.HDMI-A-1.logical.x")
// if err != nil {
// 	return
// }
// utils.NiriSend(
// 	map[string]interface{}{
// 		"Output": map[string]interface{}{
// 			"output": "HDMI-A-1",
// 			"action": map[string]interface{}{"On": nil},
// 		},
// 	},
// )
// if posX.(float64) > 0 {
// utils.RunCMD("niri msg output HDMI-A-1 position set 0 0", false)
// utils.NiriSend(
// 	map[string]interface{}{
// 		"Output": map[string]interface{}{
// 			"output": "HDMI-A-1",
// 			"action": map[string]interface{}{
// 				"Position": map[string]interface{}{
// 					"Specific": map[string]interface{}{"x": 0, "y": 0},
// 				}},
// 		},
// 	},
// )
// } else {
// utils.NiriSend(
// 	map[string]interface{}{
// 		"OutputActionCon": map[string]interface{}{
// 			"Output": "HDMI-A-1",
// 			"Action": map[string]interface{}{
// 				"Position": map[string]interface{}{
// 					"Specific": map[string]interface{}{"x": 1920, "y": 0},
// 				}},
// 		},
// 	},
// )
// utils.RunCMD("niri msg output HDMI-A-1 position set 1920 0", false)
// }
// utils.RunCMD("notify-send 'Toggle HDMI-A-1 Position'", false)
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
// }
