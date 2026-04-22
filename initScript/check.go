package initScript

import (
	"fmt"
	"niri-startup/command/action"
	"niri-startup/utils"
)

func Check(PORT int) (err error) {
	/** 检查是否有niri-startup已经运行, 如果有关闭 */
	pid, err := utils.RunCMD("pgrep niri-startup", true)
	if err != nil {
		return err
	}
	if len(pid) > 0 {
		utils.RunCMD(fmt.Sprintf("pkill %s", pid), true)
	}

	/** 检查端口是否被占用 */
	if !utils.IsPortAvailable(PORT) {
		return fmt.Errorf("端口 %d 已被占用", PORT)
	}

	// 调整家里和公司的第二个屏幕的位置
	outputsStr, err := utils.NiriSend(map[string]interface{}{
		"Outputs": nil,
	})
	if err != nil {
		return
	}
	HDMIOutputMake, err := utils.GetNestedValueFromStr(outputsStr, "Ok.Outputs.HDMI-A-1.make")
	if HDMIOutputMake == "Huawei Technologies Co., Inc." {
		action.ToggleHdmiOutoutPosition()
	}

	return
}
