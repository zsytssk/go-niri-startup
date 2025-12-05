package state

import "encoding/json"

type Msg struct {
	WorkspacesChanged      *WorkspacesChanged      `json:"WorkspacesChanged,omitempty"`
	WorkspaceActivated     *WorkspaceActivated     `json:"WorkspaceActivated,omitempty"`
	WindowsChanged         *WindowsChanged         `json:"WindowsChanged,omitempty"`
	WindowClosed           *WindowClosed           `json:"WindowClosed,omitempty"`
	OverviewOpenedOrClosed *OverviewOpenedOrClosed `json:"OverviewOpenedOrClosed,omitempty"`
	WindowOpenedOrChanged  *WindowOpenedOrChanged  `json:"WindowOpenedOrChanged,omitempty"`
	WindowFocusChanged     *WindowFocusChanged     `json:"WindowFocusChanged,omitempty"`
	WindowLayoutsChanged   *WindowLayoutsChanged   `json:"WindowLayoutsChanged,omitempty"`
	ScreenshotCaptured     *ScreenshotCaptured     `json:"ScreenshotCaptured,omitempty"`
}

type WorkspacesChanged struct {
	Workspaces []Workspace `json:"workspaces"`
}
type WindowsChanged struct {
	Windows []Window `json:"windows"`
}
type WorkspaceActivated struct {
	Id      int  `json:"id"`
	Focused bool `json:"focused"`
}
type WindowClosed struct {
	Id int `json:"id"`
}
type ScreenshotCaptured struct {
	Path string `json:"path"`
}

type OverviewOpenedOrClosed struct {
	IsOpen bool `json:"is_open"`
}

type WindowOpenedOrChanged struct {
	Window Window `json:"window"`
}
type WindowFocusChanged = WindowClosed
type WindowLayoutsChanged struct {
	Changes [][2]json.RawMessage `json:"changes"`
}

type Logical struct {
	X         int     `json:"x"`
	Y         int     `json:"y"`
	Width     int     `json:"width"`
	Height    int     `json:"height"`
	Scale     float64 `json:"scale"`
	Transform string  `json:"transform"`
}

type Monitor struct {
	Name         string        `json:"name"`
	Make         string        `json:"make"`
	Model        string        `json:"model"`
	Serial       string        `json:"serial"`
	PhysicalSize []int         `json:"physical_size"`
	Modes        []interface{} `json:"modes"` // 如果 modes 里是字符串数组
	CurrentMode  int           `json:"current_mode"`
	IsCustomMode bool          `json:"is_custom_mode"`
	VRREnabled   bool          `json:"vrr_enabled"`
	VRRSupported bool          `json:"vrr_supported"`
	Logical      Logical       `json:"logical"`
}
