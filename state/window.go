package state

type FocusTimestamp struct {
	Secs  int `json:"secs"`
	Nanos int `json:"nanos"`
}

type WindowLayout struct {
	PosInScrollingLayout   [2]int     `json:"pos_in_scrolling_layout"`
	TileSize               [2]float64 `json:"tile_size"`
	WindowSize             [2]int     `json:"window_size"`
	TilePosInWorkspaceView *[2]int    `json:"tile_pos_in_workspace_view"` // null 对应指针
	WindowOffsetInTile     [2]float64 `json:"window_offset_in_tile"`
}

type Window struct {
	ID             int            `json:"id"`
	Title          string         `json:"title"`
	AppId          string         `json:"app_id"`
	PID            int            `json:"pid"`
	WorkspaceID    int            `json:"workspace_id"`
	IsFocused      bool           `json:"is_focused"`
	IsFloating     bool           `json:"is_floating"`
	IsUrgent       bool           `json:"is_urgent"`
	Layout         WindowLayout   `json:"layout"`
	FocusTimestamp FocusTimestamp `json:"focus_timestamp"`
}
