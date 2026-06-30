package cmdAction

import (
	"niri-startup/action"
	"niri-startup/utils"
)

var HdmiOutoutPositionAction = CmdActionItem {
	CmdList: []string {"ToggleHdmiOutoutPosition"},
	Fn: ToggleHdmiOutoutPosition,
}

func ToggleHdmiOutoutPosition(_cmd string) error {
	outputsStr, err := utils.NiriSend(map[string]interface{}{
		"Outputs": nil,
	})
	if err != nil {
		return err
	}

	posX, err := utils.GetNestedValueFromStr(outputsStr, "Ok.Outputs.HDMI-A-1.logical.x")
	if err != nil {
		return err
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

	return nil
}

