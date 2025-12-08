package utils

import "time"

type Empty struct{}
type WindowWithId struct {
	Id int `json:"id"`
}
type WindowReference struct {
	Id   int    `json:"Id,omitempty"`
	Name string `json:"Name,omitempty"`
}
type MoveWindowToWorkspace struct {
	WindowId  int             `json:"window_id,omitempty"`
	Focus     bool            `json:"focus"`
	Reference WindowReference `json:"reference"`
}
type SetWindowSizeChange struct {
	SetFixed int `json:"SetFixed,omitempty"`
}
type SetWindowSize struct {
	Id     int                 `json:"id"`
	Change SetWindowSizeChange `json:"change"`
}
type ScreenshotScreen struct {
	WriteToDisk bool `json:"write_to_disk,omitempty"`
	ShowPointer bool `json:"show_pointer,omitempty"`
}
type MoveWorkspaceToMonitor struct {
	Output    string          `json:"output,omitempty"`
	Reference WindowReference `json:"reference,omitempty"`
}
type MoveWorkspaceToIndex struct {
	Index     int             `json:"index,omitempty"`
	Reference WindowReference `json:"reference,omitempty"`
}
type FocusWorkspace = MoveWorkspaceToIndex
type Action struct {
	FocusWorkspace         *FocusWorkspace         `json:"FocusWorkspace,omitempty"`
	MoveWorkspaceToMonitor *MoveWorkspaceToMonitor `json:"MoveWorkspaceToMonitor,omitempty"`
	MoveWorkspaceToIndex   *MoveWorkspaceToIndex   `json:"MoveWorkspaceToIndex,omitempty"`
	MoveWindowToWorkspace  *MoveWindowToWorkspace  `json:"MoveWindowToWorkspace,omitempty"`
	SetWindowHeight        *SetWindowSize          `json:"SetWindowHeight,omitempty"`
	SetWindowWidth         *SetWindowSize          `json:"SetWindowWidth,omitempty"`
	CenterWindow           *WindowWithId           `json:"CenterWindow,omitempty"`
	FocusWindow            *WindowWithId           `json:"FocusWindow,omitempty"`
	ToggleWindowFloating   *WindowWithId           `json:"ToggleWindowFloating,omitempty"`
	MoveWindowToFloating   *WindowWithId           `json:"MoveWindowToFloating,omitempty"`
	PowerOffMonitors       *Empty                  `json:"PowerOffMonitors,omitempty"`
	ScreenshotScreen       *ScreenshotScreen       `json:"ScreenshotScreen,omitempty"`
	Screenshot             *ScreenshotScreen       `json:"Screenshot,omitempty"`
	ScreenshotWindow       *ScreenshotScreen       `json:"ScreenshotWindow,omitempty"`
	Sleep                  int                     `json:"-"`
}

var SocketInstance *Client

func GetSocketInstance() *Client {
	if SocketInstance == nil {
		SocketInstance = NewClient("SendClient")
		SocketInstance.Connect()
	}
	return SocketInstance
}

func NiriSendAction(obj Action) {
	client := GetSocketInstance()
	var msg = map[string]Action{"Action": obj}
	client.Send(msg)
}

func NiriSendActionArr(arr []Action) {
	client := GetSocketInstance()
	for _, obj := range arr {
		if obj.Sleep != 0 {
			time.Sleep(time.Duration(obj.Sleep) * time.Millisecond)
		} else {
			var msg = map[string]Action{"Action": obj}
			client.Send(msg)
		}
	}
}
