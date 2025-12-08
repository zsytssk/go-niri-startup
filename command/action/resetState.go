package action

import (
	"niri-startup/state"
	"niri-startup/utils"
)

func ResetState() {
	instance := state.GetStateInstance()
	CurrentWindowId := instance.CurrentWindowId
	Windows := instance.Windows
	Workspaces := instance.Workspaces
	OriginWindowInfo := instance.OriginWindowInfo
	OriginWorkspaceInfo := instance.OriginWorkspaceInfo

	activeWorkspaces := make([]state.Workspace, 0)

	for _, workspace := range Workspaces {
		if workspace.IsFocused {
			activeWorkspaces = append(activeWorkspaces, workspace)
		} else if workspace.IsActive {
			activeWorkspaces = append([]state.Workspace{workspace}, activeWorkspaces...)
		}
	}

	for _, win := range Windows {
		oriInfo, ok := OriginWindowInfo[win.ID]
		if !ok {
			delete(OriginWindowInfo, win.ID)
			continue
		}
		if win.WorkspaceID != oriInfo.Workspace {
			utils.NiriSendAction(utils.Action{MoveWindowToWorkspace: &utils.MoveWindowToWorkspace{
				WindowId:  win.ID,
				Focus:     false,
				Reference: utils.WindowReference{Id: oriInfo.Workspace},
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
			utils.NiriSendActionArr([]utils.Action{
				{
					MoveWorkspaceToMonitor: &utils.MoveWorkspaceToMonitor{
						Output: oriInfo.Outout,
						Reference: utils.WindowReference{
							Id: workspace.ID,
						},
					},
				},
				{
					MoveWorkspaceToIndex: &utils.MoveWorkspaceToIndex{
						Index: oriInfo.Idx,
						Reference: utils.WindowReference{
							Id: workspace.ID,
						},
					},
				},
			})
		}
	}

	for _, workspace := range activeWorkspaces {
		utils.NiriSendAction(utils.Action{
			FocusWorkspace: &utils.FocusWorkspace{
				Reference: utils.WindowReference{Id: workspace.ID},
			}},
		)
	}

	utils.NiriSendAction(utils.Action{
		FocusWindow: &utils.WindowWithId{
			Id: CurrentWindowId,
		}},
	)
}
