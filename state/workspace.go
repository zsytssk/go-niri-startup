package state

type Workspace struct {
	ID             int     `json:"id"`
	Idx            int     `json:"idx"`
	Name           *string `json:"name"` // null 用指针接收
	Output         string  `json:"output"`
	IsUrgent       bool    `json:"is_urgent"`
	IsActive       bool    `json:"is_active"`
	IsFocused      bool    `json:"is_focused"`
	ActiveWindowID int     `json:"active_window_id"`
}
