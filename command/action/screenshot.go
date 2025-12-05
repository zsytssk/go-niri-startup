package action

import (
	"fmt"
	"niri-startup/state"
	"niri-startup/utils"
)

func Screenshot() {
	instance := state.GetStateInstance()
	waitScreenShot := state.UseWaitScreenShot(instance)
	path := waitScreenShot()
	if path == "" {
		return
	}
	utils.RunCMD(
		fmt.Sprintf(`satty --filename %s --actions-on-enter save-to-clipboard`, path),
		false,
	)
}
