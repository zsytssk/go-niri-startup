package state

import (
	"niri-startup/utils"
	"sync"
)

var (
	StateInstance State
	once          sync.Once
)

func GetStateInstance() *State {

	once.Do(func() {
		create()
	})
	return &StateInstance

}

func create() {

	StateInstance = State{
		Windows:             make(map[int]*Window, 0),
		Workspaces:          make(map[int]*Workspace, 0),
		OriginWorkspaceInfo: make(map[int]*OriginWorkspaceInfo, 0),
		OriginWindowInfo:    make(map[int]*OriginWindowInfo, 0),
		Event: utils.Event{
			Listeners: make(map[string][]utils.Listener, 0),
			Counter:   0,
		},
		Client: utils.NewClient("EventStreamClient"),
	}
	StateInstance.Client.Connect()
	StateInstance.BindEventStream()
}
