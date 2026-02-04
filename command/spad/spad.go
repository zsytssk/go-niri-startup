package spad

import (
	"encoding/json"
	"log"
	"net/http"
	"niri-startup/action"
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
	view := state.GetStateSnapshot()
	waitWindowOpen := state.UseWaitWindowOpen(instance)
	windowFilter := state.UseWindowFilter(instance)
	onWindowBlur := state.UseOnWindowBlur(instance)
	updateSpadOriInfo := state.UseUpdateSpadOriInfo(instance)
	matchFn := UseMatchFn(spadConf)
	currentWorkspaceId := view.CurrentWorkspaceId

	var win *state.Window
	wins := windowFilter(matchFn)
	if len(wins) > 0 {
		win = wins[0]
		if off, ok := BindFnMap[win.ID]; ok {
			off()
			delete(BindFnMap, win.ID)
		}
		if isSpadWinAction(win) {
			utils.NiriSendActionArr([]action.Action{
				{
					MoveWindowToWorkspace: &action.MoveWindowToWorkspace{
						WindowId:  win.ID,
						Focus:     false,
						Reference: action.WindowReference{Name: config.SpadWorkspaceName},
					},
				},
				{
					ToggleWindowFloating: &action.WindowWithId{
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
		if err != nil {
			log.Println("can't find open window for spad", req.Name)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		updateSpadOriInfo(win.ID)
	}

	utils.NiriSendActionArr([]action.Action{
		{
			MoveWindowToWorkspace: &action.MoveWindowToWorkspace{
				WindowId:  win.ID,
				Focus:     false,
				Reference: action.WindowReference{Id: currentWorkspaceId},
			},
		},
		{
			SetWindowHeight: &action.SetWindowSize{
				Id: win.ID, Change: action.SetWindowSizeChange{SetFixed: spadConf.Height},
			},
		},
		{
			SetWindowWidth: &action.SetWindowSize{
				Id: win.ID, Change: action.SetWindowSizeChange{SetFixed: spadConf.Width},
			},
		},
		{
			MoveWindowToFloating: &action.WindowWithId{Id: win.ID},
		},
		{
			FocusWindow: &action.WindowWithId{Id: win.ID},
		},
		{Sleep: 80},
		{
			CenterWindow: &action.WindowWithId{Id: win.ID},
		},
	})
	win.WorkspaceID = currentWorkspaceId

	BindFnMap[win.ID] = onWindowBlur(win, func() {
		delete(BindFnMap, win.ID)
		utils.NiriSendActionArr([]action.Action{
			{
				MoveWindowToWorkspace: &action.MoveWindowToWorkspace{
					WindowId:  win.ID,
					Focus:     false,
					Reference: action.WindowReference{Name: config.SpadWorkspaceName},
				},
			},
			{
				ToggleWindowFloating: &action.WindowWithId{
					Id: win.ID,
				},
			},
		})
	})
}
