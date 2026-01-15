package action

import (
	"niri-startup/config"
	"niri-startup/state"
	"niri-startup/utils"
)

var matchFn = func(w *state.Window) bool {
	if w.Title == "画中画" {
		return true
	}
	return false
}

func togglePictureInPicture() {
	instance := state.GetStateInstance()
	windowFilter := state.UseWindowFilter(instance)
	wins := windowFilter(matchFn)
	currentWorkspaceId := instance.CurrentWorkspaceId

	if len(wins) == 0 {
		return
	}
	win := wins[0]
	if win.WorkspaceID == currentWorkspaceId {
		utils.NiriSendActionArr([]utils.Action{
			{
				MoveWindowToWorkspace: &utils.MoveWindowToWorkspace{
					WindowId:  win.ID,
					Focus:     false,
					Reference: utils.WindowReference{Name: config.SpadWorkspaceName},
				},
			},
		})
	} else {
		utils.NiriSendActionArr([]utils.Action{
			{
				MoveWindowToWorkspace: &utils.MoveWindowToWorkspace{
					WindowId:  win.ID,
					Focus:     false,
					Reference: utils.WindowReference{Id: currentWorkspaceId},
				},
			},
			// {
			// 	FocusWindow: &utils.WindowWithId{Id: win.ID},
			// },
		})
	}
}
