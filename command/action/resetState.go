package action

import (
	"niri-startup/action"
	"niri-startup/state"
	"niri-startup/utils"
)

func ResetState() {
	view := state.GetStateSnapshot()
	CurrentWindowId := view.CurrentWindowId
	Windows := view.Windows
	Workspaces := view.Workspaces
	OriginWindowInfo := view.OriginWindowInfo
	OriginWorkspaceInfo := view.OriginWorkspaceInfo

	activeWorkspaces := make([]*state.Workspace, 0)

	for _, workspace := range Workspaces {
		if workspace.IsFocused {
			activeWorkspaces = append(activeWorkspaces, workspace)
		} else if workspace.IsActive {
			activeWorkspaces = append([]*state.Workspace{workspace}, activeWorkspaces...)
		}
	}

	for _, win := range Windows {
		oriInfo, ok := OriginWindowInfo[win.ID]
		if !ok {
			delete(OriginWindowInfo, win.ID)
			continue
		}
		if win.WorkspaceID != oriInfo.Workspace {
			utils.NiriSendAction(action.Action{
				MoveWindowToWorkspace: &action.MoveWindowToWorkspace{
					WindowId:  win.ID,
					Focus:     false,
					Reference: action.WindowReference{Id: oriInfo.Workspace},
				}})
		}
	}
	for _, workspace := range Workspaces {
		oriInfo, ok := OriginWorkspaceInfo[workspace.ID]
		if !ok {
			delete(OriginWorkspaceInfo, workspace.ID)
			continue
		}
		if workspace.Output != oriInfo.Outout {
			utils.NiriSendActionArr([]action.Action{
				{
					MoveWorkspaceToMonitor: &action.MoveWorkspaceToMonitor{
						Output: oriInfo.Outout,
						Reference: action.WindowReference{
							Id: workspace.ID,
						},
					},
				},
				{
					MoveWorkspaceToIndex: &action.MoveWorkspaceToIndex{
						Index: oriInfo.Idx,
						Reference: action.WindowReference{
							Id: workspace.ID,
						},
					},
				},
			})
		}
	}

	for _, workspace := range activeWorkspaces {
		utils.NiriSendAction(action.Action{
			FocusWorkspace: &action.FocusWorkspace{
				Reference: action.WindowReference{Id: workspace.ID},
			}},
		)
	}

	utils.NiriSendAction(action.Action{
		FocusWindow: &action.WindowWithId{
			Id: CurrentWindowId,
		}},
	)
}
