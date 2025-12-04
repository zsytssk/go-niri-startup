package state

import (
	"encoding/json"
	"niri-startup/utils"
)

type State struct {
	Windows    map[int]Window
	Workspaces map[int]Workspace
	utils.Event
}

func (state *State) BindEventStream(client *utils.Client) {
	go func() {
		<-client.Connected
		client.Send("\"EventStream\"")
		for msg := range client.Message {
			key, data := utils.GetData(msg)
			switch key {
			case "WorkspacesChanged":
				{
					_, data := utils.GetData(data)
					workspaces := make([]Workspace, 0)
					json.Unmarshal(data, &workspaces)
					for _, workspace := range workspaces {
						state.Workspaces[workspace.ID] = workspace
					}
					state.TriggerEvent("WorkspacesChanged", state.Workspaces)
				}
			case "WorkspaceActivated":
				{
					_, data := utils.GetData(data)
					workspaces := make([]Workspace, 0)
					json.Unmarshal(data, &workspaces)
					for _, workspace := range workspaces {
						state.Workspaces[workspace.ID] = workspace
					}
					state.TriggerEvent("WorkspacesChanged", state.Workspaces)
				}
			case "WindowsChanged":
				{
					_, data := utils.GetData(data)
					windows := make([]Window, 0)
					json.Unmarshal(data, &windows)
					for _, window := range windows {
						state.Windows[window.ID] = window
					}

					state.TriggerEvent("WindowsChanged", state.Windows)
				}
			}
		}
	}()
}
