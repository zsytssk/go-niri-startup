package action

import (
	"fmt"
	"niri-startup/utils"
	"strings"
)

func PickColor() {
	str, err := utils.RunCMD("niri msg pick-color", false)
	if err != nil {
		return
	}
	lines := strings.Split(str, "\n")
	if len(lines) < 2 {
		return
	}
	color := strings.Split(lines[1], ": ")[1]
	utils.RunCMD(fmt.Sprintf(`echo -n "%s" | wl-copy`, color), false)
}
