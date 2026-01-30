package action

import (
	"niri-startup/state"
	"niri-startup/utils"
	"slices"
	"sync/atomic"
)

var isSwitch atomic.Bool

func SwitchScreen(changeSpace int) {
	// log.Println(`test:>SwitchScreen:>1:>isSwitch`, isSwitch.Load())
	instance := state.GetStateInstance()
	waitWorkspaceChangeComplete := state.UseWorkspaceChangeComplete(instance)
	view := instance.GetSnapshot()
	workspaces := view.Workspaces
	curWorkspace, ok := workspaces[view.CurrentWorkspaceId]
	if !ok || isSwitch.Load() {
		return
	}
	isSwitch.Store(true)
	curOutput := curWorkspace.Output
	curIndex := slices.Index(view.Outputs, curOutput)
	nextIndex := curIndex + changeSpace

	if nextIndex >= len(view.Outputs) {
		nextIndex = 0
	} else if nextIndex < 0 {
		nextIndex = len(view.Outputs) - 1
	}
	nextOutput := view.Outputs[nextIndex]
	if !ok {
		return
	}

	curOutputWorkspaces := make([]*state.Workspace, 0)
	nextOutputWorkspaces := make([]*state.Workspace, 0)

	for _, item := range workspaces {
		switch item.Output {
		case curOutput:
			curOutputWorkspaces = append(curOutputWorkspaces, item)
		case nextOutput:
			nextOutputWorkspaces = append(nextOutputWorkspaces, item)
		}
	}

	slices.SortFunc(curOutputWorkspaces, func(a, b *state.Workspace) int {
		return a.Idx - b.Idx
	})
	slices.SortFunc(nextOutputWorkspaces, func(a, b *state.Workspace) int {
		return a.Idx - b.Idx
	})

	moveActions := []utils.Action{}
	focusActions := []utils.Action{}
	changeList := make([]state.ChangeWorkspaceInfo, 0)
	for _, workspace := range append(nextOutputWorkspaces, curOutputWorkspaces...) {
		var goOutput string
		if workspace.Output == curOutput {
			goOutput = nextOutput
		} else {
			goOutput = curOutput
		}
		changeList = append(changeList, state.ChangeWorkspaceInfo{
			ID:       workspace.ID,
			Output:   goOutput,
			Idx:      workspace.Idx,
			IsActive: workspace.IsActive,
		})
		actions := []utils.Action{
			{
				MoveWorkspaceToMonitor: &utils.MoveWorkspaceToMonitor{
					Output: goOutput,
					Reference: utils.WindowReference{
						Id: workspace.ID,
					},
				},
			},
			{
				MoveWorkspaceToIndex: &utils.MoveWorkspaceToIndex{
					Index: workspace.Idx,
					Reference: utils.WindowReference{
						Id: workspace.ID,
					},
				},
			},
		}
		if workspace.IsFocused {
			focusActions = append(focusActions, utils.Action{
				FocusWorkspace: &utils.FocusWorkspace{
					Reference: utils.WindowReference{
						Id: workspace.ID,
					},
				}},
			)
		} else if workspace.IsActive {
			actions = append([]utils.Action{
				{FocusWorkspace: &utils.FocusWorkspace{
					Reference: utils.WindowReference{
						Id: workspace.ID,
					},
				}},
			}, actions...)
		}
		moveActions = append(moveActions, actions...)
	}
	changeCompleteCh := waitWorkspaceChangeComplete(changeList)
	utils.NiriSendActionArr(moveActions)
	utils.NiriSendActionArr(focusActions)

	// wait switch animation end
	<-changeCompleteCh
	isSwitch.Store(false)
	// log.Println(`test:>SwitchScreen:>2:>isSwitch`, isSwitch.Load())
}
