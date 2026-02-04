package action

import (
	"niri-startup/action"
	"niri-startup/state"
	"niri-startup/utils"
	"slices"
)

func FocusNextWindow() {
	instance := state.GetStateInstance()
	view := instance.GetSnapshot()
	currentWorkspaceId := view.CurrentWorkspaceId
	getWindows := state.UseWorkspaceWindows(instance)

	windows := getWindows(currentWorkspaceId)
	if len(windows) == 0 {
		return
	}
	curIndex := slices.IndexFunc(windows, func(w *state.Window) bool {
		return w.IsFocused
	})
	nextIndex := curIndex + 1
	if nextIndex == len(windows) {
		nextIndex = 0
	}
	nextWindow := windows[nextIndex]

	utils.NiriSendAction(action.Action{
		FocusWindow: &action.WindowWithId{Id: nextWindow.ID},
	})
}
