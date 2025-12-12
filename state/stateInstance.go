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
		Event:               utils.NewEvent(),
		Client:              utils.NewClient("EventStreamClient"),
	}
	StateInstance.Client.Connect()
	StateInstance.BindEventStream()
}
