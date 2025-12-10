package state

import (
	"encoding/json"
	"niri-startup/utils"
	"sort"
)

type OriginWorkspaceInfo struct {
	Outout string
	Idx    int
}
type OriginWindowInfo struct {
	Workspace int
}
type State struct {
	Outputs             []string
	CurrentWindowId     int
	CurrentWorkspaceId  int
	OverviewOpen        bool
	Windows             map[int]Window
	Workspaces          map[int]Workspace
	OriginWorkspaceInfo map[int]OriginWorkspaceInfo
	OriginWindowInfo    map[int]OriginWindowInfo
	utils.Event         `json:"-"`
}

func (s *State) BindEventStream(client *utils.Client) {
	go func() {
		<-client.Connected
		client.Send("EventStream")
		for msg := range client.ReviveMsgCh {
			msgType := utils.GetMsgType(msg)
			var data Msg
			json.Unmarshal(msg, &data)
			s.TriggerEvent(msgType, data)
			// fmt.Println(`test:>msg`, msgType, string(msg))
			switch msgType {
			case "WorkspacesChanged":
				{
					s.Workspaces = make(map[int]Workspace, 0)
					s.outputsChange(data.WorkspacesChanged.Workspaces)
					for _, m := range data.WorkspacesChanged.Workspaces {
						s.addWorkspace(m)
					}
				}
			case "WorkspaceActivated":
				{
					w := data.WorkspaceActivated
					s.setActiveWorkspace(
						w.Id,
						0,
						w.Focused,
					)
				}
			case "WindowsChanged":
				{
					s.Windows = make(map[int]Window, 0)
					wins := data.WindowsChanged.Windows
					for _, wi := range wins {
						s.addWindow(wi)
					}
				}
			case "WindowClosed":
				{
					s.windowClose(data.WindowClosed.Id)
				}
			case "OverviewOpenedOrClosed":
				{
					s.OverviewOpen = data.OverviewOpenedOrClosed.IsOpen
				}
			case "WindowOpenedOrChanged":
				{
					s.addWindow(data.WindowOpenedOrChanged.Window)
				}
			case "WindowFocusChanged":
				{
					s.setCurWindowId(data.WindowFocusChanged.Id)
				}
			case "WindowLayoutsChanged":
				{
					for _, w := range data.WindowLayoutsChanged.Changes {

						var id int
						json.Unmarshal(w[0], &id)
						var layout WindowLayout
						json.Unmarshal(w[1], &layout)

						if w, ok := s.Windows[id]; ok {
							w.Layout = layout
							s.Windows[id] = w
						}
					}
				}
			}
		}
	}()
}

func (s *State) setActiveWorkspace(curId int, activeWindowId int, focus bool) {
	var output string
	if w, ok := s.Workspaces[curId]; ok {
		output = w.Output
	} else {
		return
	}

	for id, item := range s.Workspaces {
		if item.Output != output {
			continue
		}
		item.IsActive = item.ID == curId
		if activeWindowId != 0 {
			if item.IsActive {
				item.ActiveWindowID = activeWindowId
			} else {
				item.ActiveWindowID = 0
			}
		}
		s.Workspaces[id] = item
	}

	if focus {
		s.CurrentWorkspaceId = curId

		for _, item := range s.Workspaces {
			item.IsFocused = item.ID == curId
			s.Workspaces[item.ID] = item
		}
	}
}
func (s *State) GetWindowOutput(windowId int) string {
	var window Window
	if w, ok := s.Windows[windowId]; ok {
		window = w
	} else {
		return ""
	}

	if w, ok := s.Workspaces[window.WorkspaceID]; ok {
		return w.Output
	} else {
		return ""
	}

}

func (s *State) addWorkspace(workspace Workspace) {
	s.Workspaces[workspace.ID] = workspace
	if workspace.IsFocused {
		s.setActiveWorkspace(
			workspace.ID,
			0,
			true,
		)
	}

	if _, ok := s.OriginWorkspaceInfo[workspace.ID]; !ok {
		s.OriginWorkspaceInfo[workspace.ID] = OriginWorkspaceInfo{
			Outout: workspace.Output,
			Idx:    workspace.Idx,
		}
	}
}

func (s *State) setCurWindowId(curId int) {
	s.CurrentWindowId = curId

	w, ok := s.Windows[curId]
	if ok {
		for _, item := range s.Windows {
			item.IsFocused = item.ID == curId
			s.Windows[item.ID] = item
		}
		workspaceId := w.WorkspaceID
		if workspaceId != 0 {
			s.setActiveWorkspace(workspaceId, curId, true)
		}
	}
	s.TriggerEvent("FocusWindow", w)
}

func (s *State) windowClose(id int) {

	w, ok := s.Windows[id]
	if !ok {
		return
	}
	delete(s.Windows, id)
	if workspace, ok := s.Workspaces[w.WorkspaceID]; ok && workspace.ActiveWindowID == id {
		workspace.ActiveWindowID = 0
	}
}
func (s *State) addWindow(window Window) {
	s.Windows[window.ID] = window
	if window.IsFocused {
		// @todo setTimeout
		// s.setCurWindowId(window.ID)
	}

	if _, ok := s.OriginWindowInfo[window.ID]; !ok {
		s.OriginWindowInfo[window.ID] = OriginWindowInfo{
			Workspace: window.WorkspaceID,
		}
	}
}
func (s *State) outputsChange(workspaces []Workspace) {
	var newOutputs = make(map[string]bool)

	for _, workspace := range workspaces {
		newOutputs[workspace.Output] = true
	}
	if len(newOutputs) == len(s.Outputs) {
		var isSame = true
		for _, o := range s.Outputs {
			if _, ok := newOutputs[o]; !ok {
				isSame = false
				break
			}
		}
		if isSame {
			return
		}
	}

	outputsStr, err := utils.RunCMD("niri msg --json outputs", false)
	if err != nil {
		return
	}
	var monitors map[string]Monitor
	if err := json.Unmarshal([]byte(outputsStr), &monitors); err != nil {
		panic(err)
	}

	var arr = make([]Monitor, 0, len(monitors))
	for _, m := range monitors {
		arr = append(arr, m)
	}
	sort.Slice(arr, func(i, j int) bool {
		a := arr[i]
		b := arr[j]
		return a.Logical.X < b.Logical.X || a.Logical.Y < b.Logical.Y
	})
	s.Outputs = s.Outputs[:0]

	for _, m := range arr {
		s.Outputs = append(s.Outputs, m.Name)
	}
}
