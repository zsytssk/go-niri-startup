package action

import (
	"fmt"
	"niri-startup/state"
	"niri-startup/utils"
	"slices"
	"sync"
)

var isSwitch = false

func SwitchScreen(changeSpace int) {
	instance := state.GetStateInstance()
	workspaces := instance.Workspaces
	curWorkspace, ok := workspaces[instance.CurrentWorkspaceId]
	if !ok || isSwitch {
		return
	}
	isSwitch = true
	curOutput := curWorkspace.Output
	curIndex := slices.Index(instance.Outputs, curOutput)
	nextIndex := curIndex + changeSpace

	if nextIndex > len(instance.Outputs) {
		nextIndex = 0
	} else if nextIndex < 0 {
		nextIndex = len(instance.Outputs) - 1
	}
	nextOutput := instance.Outputs[nextIndex]
	if !ok {
		return
	}

	curOutputWorkspaces := make([]state.Workspace, 0)
	nextOutputWorkspaces := make([]state.Workspace, 0)

	for _, item := range workspaces {
		switch item.Output {
		case curOutput:
			curOutputWorkspaces = append(curOutputWorkspaces, item)
		case nextOutput:
			nextOutputWorkspaces = append(nextOutputWorkspaces, item)
		}
	}

	slices.SortFunc(curOutputWorkspaces, func(a, b state.Workspace) int {
		return a.Idx - b.Idx
	})
	slices.SortFunc(nextOutputWorkspaces, func(a, b state.Workspace) int {
		return a.Idx - b.Idx
	})
	var wg sync.WaitGroup
	num := len(curOutputWorkspaces) + len(nextOutputWorkspaces)
	// chArr := make([]chan int, num)
	wg.Add(num)

	for _, workspace := range append(curOutputWorkspaces, nextOutputWorkspaces...) {

		go func() {
			var goOutput string
			if workspace.Output == curOutput {
				goOutput = nextOutput
			} else {
				goOutput = curOutput
			}
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
				actions = append(actions, utils.Action{FocusWorkspace: &utils.FocusWorkspace{
					Reference: utils.WindowReference{
						Id: workspace.ID,
					},
				}})
			} else if workspace.IsActive {
				actions = append([]utils.Action{
					{FocusWorkspace: &utils.FocusWorkspace{
						Reference: utils.WindowReference{
							Id: workspace.ID,
						},
					}},
				}, actions...)
			}
			utils.NiriSendActionArr(actions)
			// <-chArr[index]
			defer wg.Done()

		}()
	}
	fmt.Println(`wg:>1`)
	wg.Wait()
	fmt.Println(`wg:>2`)
	isSwitch = false
}
