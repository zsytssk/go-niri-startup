package initScript

import (
	"fmt"
	"niri-startup/utils"
)

/** 初始运行脚本 */
func Run() {
	outputsStr, err := utils.NiriSend(map[string]interface{}{
		"Outputs": nil,
	})
	if err != nil {
		return
	}
	HDMIOutputMake, err := utils.GetNestedValueFromStr(outputsStr, "Ok.Outputs.HDMI-A-1.make")
	HDMIOutputModel, err := utils.GetNestedValueFromStr(outputsStr, "Ok.Outputs.HDMI-A-1.model")
	utils.RunCMD(fmt.Sprintf("notify-send HDMIOutputMake \"%s\"", HDMIOutputMake), false)
	utils.RunCMD(fmt.Sprintf("notify-send HDMIOutputModel \"%s\"", HDMIOutputModel), false)
}
