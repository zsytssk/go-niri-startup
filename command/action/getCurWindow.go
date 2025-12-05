package action

import (
	"fmt"
	"net/http"
	"niri-startup/state"
	"niri-startup/utils"
	"slices"
)

func GetCurWindow(write http.ResponseWriter, req ActionReq) {
	bindActiveWindowChange(req)
	instance := state.GetStateInstance()
	getWindows := state.UseWorkspaceWindows(instance)
	windows := instance.Windows
	workspaces := instance.Workspaces

	var activeWindowId int
	var workspaceWindows []*state.Window
	for _, workspace := range workspaces {
		if workspace.Output == req.Output && workspace.IsActive {
			activeWindowId = workspace.ActiveWindowID
			workspaceWindows = getWindows(workspace.ID)
			break
		}
	}
	if activeWindowId == 0 {
		utils.ReturnHttp(write, "")
		return
	}
	index := slices.IndexFunc(workspaceWindows, func(w *state.Window) bool {
		return w.ID == activeWindowId
	})
	allNum := len(workspaceWindows)
	cur, ok := windows[activeWindowId]
	if !ok {
		utils.ReturnHttp(write, "")
		return
	}
	// fmt.Println(`test:>`, req.Output, activeWindowId, instance.CurrentWindowId)
	utils.ReturnHttp(write, fmt.Sprintf(`%d/%d %s`, index+1, allNum, cur.Title))
}

var bindMap = make(map[string]bool)

func bindActiveWindowChange(req ActionReq) {
	_, ok := bindMap[req.Output]
	if ok {
		return
	}
	bindMap[req.Output] = true
	instance := state.GetStateInstance()
	instance.OnEvent("FocusWindow", func(i interface{}) {
		w := i.(state.Window)
		if w.ID == 0 {
			utils.RunCMD(fmt.Sprintf(`pkill -RTMIN+%d waybar`, req.Signal), false)
			return
		}
		_, ok := instance.Workspaces[w.WorkspaceID]
		if !ok {
			return
		}
		utils.RunCMD(fmt.Sprintf(`pkill -RTMIN+%d waybar`, req.Signal), false)
	})
}
