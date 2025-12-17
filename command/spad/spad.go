package spad

import (
	"encoding/json"
	"net/http"
	"niri-startup/config"
	"niri-startup/state"
	"niri-startup/utils"
)

type SpadReq struct {
	Name string `json:"name"`
}

var BindFnMap = make(map[int]func(), 0)

func Spad(w http.ResponseWriter, r *http.Request) {
	var req SpadReq

	// 解析 JSON
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	spadConf, err := config.GetSpadConfig(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	instance := state.GetStateInstance()
	waitWindowOpen := state.UseWaitWindowOpen(instance)
	windowFilter := state.UseWindowFilter(instance)
	onWindowBlur := state.UseOnWindowBlur(instance)
	updateSpadOriInfo := state.UseUpdateSpadOriInfo(instance)
	matchFn := UseMatchFn(spadConf)
	currentWorkspaceId := instance.CurrentWorkspaceId

	var win *state.Window
	wins := windowFilter(matchFn)
	if len(wins) > 0 {
		win = wins[0]
		if off, ok := BindFnMap[win.ID]; ok {
			off()
			delete(BindFnMap, win.ID)
		}
		if isSpadWinAction(win) {
			utils.NiriSendActionArr([]utils.Action{
				{
					MoveWindowToWorkspace: &utils.MoveWindowToWorkspace{
						WindowId:  win.ID,
						Focus:     false,
						Reference: utils.WindowReference{Name: config.SpadWorkspaceName},
					},
				},
				{
					ToggleWindowFloating: &utils.WindowWithId{
						Id: win.ID,
					},
				},
			})
			return
		}
	} else {
		_, err := utils.RunCMD(spadConf.Cmd, true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		win, err = waitWindowOpen(matchFn)
		updateSpadOriInfo(win.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// fmt.Println(`test:>`, currentWorkspaceId, win.Title)
	// if strings.Contains(win.Title, "DeepSeek") {
	// 	return
	// }
	utils.NiriSendActionArr([]utils.Action{
		{
			MoveWindowToWorkspace: &utils.MoveWindowToWorkspace{
				WindowId:  win.ID,
				Focus:     false,
				Reference: utils.WindowReference{Id: currentWorkspaceId},
			},
		},
		{
			SetWindowHeight: &utils.SetWindowSize{
				Id: win.ID, Change: utils.SetWindowSizeChange{SetFixed: spadConf.Height},
			},
		},
		{
			SetWindowWidth: &utils.SetWindowSize{
				Id: win.ID, Change: utils.SetWindowSizeChange{SetFixed: spadConf.Width},
			},
		},
		{
			MoveWindowToFloating: &utils.WindowWithId{Id: win.ID},
		},
		{
			FocusWindow: &utils.WindowWithId{Id: win.ID},
		},
		{Sleep: 80},
		{
			CenterWindow: &utils.WindowWithId{Id: win.ID},
		},
	})
	win.WorkspaceID = currentWorkspaceId

	BindFnMap[win.ID] = onWindowBlur(win, func() {
		delete(BindFnMap, win.ID)
		utils.NiriSendActionArr([]utils.Action{
			{
				MoveWindowToWorkspace: &utils.MoveWindowToWorkspace{
					WindowId:  win.ID,
					Focus:     false,
					Reference: utils.WindowReference{Name: config.SpadWorkspaceName},
				},
			},
			{
				ToggleWindowFloating: &utils.WindowWithId{
					Id: win.ID,
				},
			},
		})
	})
}
