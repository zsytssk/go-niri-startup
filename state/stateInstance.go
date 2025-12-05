package state

import (
	"niri-startup/utils"
)

var StateInstance *State

func GetStateInstance() *State {

	if StateInstance == nil {
		create()
	}
	return StateInstance
}

func create() {

	client := utils.NewClient("EventStreamClient")
	client.Connect()

	var event = utils.Event{
		Listeners: make(map[string][]utils.Listener, 0),
		Counter:   0,
	}
	StateInstance = &State{
		Windows:             make(map[int]Window, 0),
		Workspaces:          make(map[int]Workspace, 0),
		OriginWorkspaceInfo: make(map[int]OriginWorkspaceInfo, 0),
		OriginWindowInfo:    make(map[int]OriginWindowInfo, 0),
		Event:               event,
	}
	StateInstance.BindEventStream(client)

	// s.OnEvent("WorkspacesChanged", func(w interface{}) {
	// 	// ws := w.(map[int]Workspace)
	// 	// fmt.Println(len(ws))
	// 	_, data := utils.GetData(w.(json.RawMessage))
	// 	workspaces := make([]Workspace, 0)
	// 	json.Unmarshal(data, &workspaces)
	// 	fmt.Println(len(workspaces))
	// })
}
