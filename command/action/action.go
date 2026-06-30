package action

import (
	"encoding/json"
	"net/http"
	"niri-startup/action"
	"niri-startup/command/action/cmdAction"
	"niri-startup/utils"
)

type ActionReq struct {
	Name   string `json:"name"`
	Output string `json:"output"`
	Signal int    `json:"signal"`
}

func Action(w http.ResponseWriter, r *http.Request) {
	var req ActionReq

	// 解析 JSON
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// var msg string
	switch req.Name {
	case "get-cur-window":
		GetCurWindow(w, req)
	case "system-actions":
		err = cmdAction.Run("System")
	case "cmd-actions":
		err = cmdAction.Run("")
	case "next-window":
		FocusNextWindow()
	case "select-window":
		SelectWindow()
	case "screenshot-screen":
		// utils.NiriSendAction(utils.Action{
		// 	ScreenshotScreen: &utils.ScreenshotScreen{WriteToDisk: true, ShowPointer: false},
		// })
		utils.RunCMD("niri msg action screenshot-screen", false)
		Screenshot()
	case "screenshot":
		// utils.NiriSendAction(utils.Action{
		// 	Screenshot: &utils.ScreenshotScreen{WriteToDisk: true, ShowPointer: false},
		// })
		utils.RunCMD("niri msg action screenshot", false)
		Screenshot()
	case "screenshot-window":
		utils.NiriSendAction(action.Action{
			ScreenshotWindow: &action.ScreenshotScreen{WriteToDisk: true},
		})
		Screenshot()
	case "toggle-input":
		ToggleInput()
	case "switch-screen-prev":
		SwitchScreen(-1)
	case "switch-screen-next":
		SwitchScreen(1)
	case "pick-color":
		PickColor()
	case "toggle-picture-in-picture":
		togglePictureInPicture()
	case "reset-state":
		ResetState()
	case "test-action":
		testAction()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
