package command

import (
	"encoding/json"
	"net/http"
	"niri-startup/state"
	"niri-startup/utils"
	"slices"
)

type RunAppReq struct {
	Cmd   string `json:"cmd"`
	AppId string `json:"app_id"`
	Title string `json:"title"`
}

func RunApp(w http.ResponseWriter, r *http.Request) {
	var req RunAppReq
	instance := state.GetStateInstance()
	waitWindowOpen := state.UseWaitWindowOpen(instance)
	windowFilter := state.UseWindowFilter(instance)

	filterFn := func(wi *state.Window) bool {
		if req.AppId != "" {
			return wi.AppId == req.AppId
		}
		if req.Title != "" {
			return wi.Title == req.Title
		}
		return false
	}
	// 解析 JSON
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	apps := windowFilter(filterFn)
	var window *state.Window
	if len(apps) != 0 {
		index := slices.IndexFunc(apps, func(item *state.Window) bool {
			return item.IsFocused
		})
		nextIndex := index + 1
		if nextIndex >= len(apps) {
			nextIndex = 0
		}
		window = apps[nextIndex]
	} else {
		_, err := utils.RunCMD(req.Cmd, true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		window, err = waitWindowOpen(filterFn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	// // CurrentWorkspaceId := instance.CurrentWorkspaceId
	// actions := []utils.Action{
	// 	{
	// 		FocusWindow: &utils.WindowWithId{
	// 			Id: window.ID,
	// 		},
	// 	},
	// 	// {
	// 	// 	CenterWindow: &utils.WindowWithId{
	// 	// 		Id: window.ID,
	// 	// 	},
	// 	// },
	// }

	utils.NiriSendAction(utils.Action{
		FocusWindow: &utils.WindowWithId{
			Id: window.ID,
		},
	})
	utils.ReturnHttp(w, "")
}
