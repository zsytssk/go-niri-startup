package state

type WorkspaceActivated struct {
	Id      int  `json:"id"`
	Focused bool `json:"focused"`
}
type WindowClosed struct {
	Id int `json:"id"`
}

type OverviewOpenedOrClosed struct {
	OverviewOpenedOrClosed bool `json:"OverviewOpenedOrClosed"`
}

type WindowOpenedOrChanged = Window
type WindowLayoutsChanged struct {
	Changes [][2]interface{} `json:"changes"`
}

type WindowLayout struct {
	PosInScrollingLayout   [2]int `json:"pos_in_scrolling_layout"`
	TileSize               [2]int `json:"tile_size"`
	WindowSize             [2]int `json:"window_size"`
	TilePosInWorkspaceView *[]int `json:"tile_pos_in_workspace_view"` // null 可表示为 nil
	WindowOffsetInTile     [2]int `json:"window_offset_in_tile"`
}
