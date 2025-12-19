package state

import (
	"fmt"
	"log"
	"niri-startup/config"
	"slices"
	"time"
)

func UseWaitWindowOpen(state *State) func(func(*Window) bool) (*Window, error) {
	return func(filterFn func(*Window) bool) (*Window, error) {
		for _, w := range state.Windows {
			if filterFn(w) {
				return w, nil
			}
		}

		ch := make(chan *Window, 1)
		var off func()
		off = state.OnEvent("WindowOpenedOrChanged", func(msg interface{}) {
			w := &msg.(EventStreamMsg).WindowOpenedOrChanged.Window
			if filterFn(w) {
				ch <- w
				off()
			}
		})

		select {
		case w := <-ch:
			return w, nil
		case <-time.After(3 * time.Second):
			off() // 超时也取消事件监听
			return nil, fmt.Errorf("timeout waiting for window")
		}
	}
}

func UseWindowFilter(state *State) func(func(*Window) bool) []*Window {
	return func(filterFn func(*Window) bool) []*Window {
		wins := make([]*Window, 0)
		for _, w := range state.Windows {
			if filterFn(w) {
				wins = append(wins, w)
			}
		}
		return wins
	}
}
func UseOnWindowBlur(state *State) func(*Window, func()) func() {
	return func(win *Window, fn func()) func() {
		var off func()
		off = state.OnEvent("FocusWindow", func(obj interface{}) {
			if obj.(*Window).ID != win.ID {
				fn()
				off()
			}
		})

		return off
	}
}

func UseWaitScreenShot(state *State) func() string {
	return func() string {
		ch := make(chan string, 1)
		var off1 func()
		var off2 func()
		off1 = state.OnEvent("ScreenshotCaptured", func(obj interface{}) {
			path := obj.(EventStreamMsg).ScreenshotCaptured.Path
			ch <- path
			off1()
			off2()
		})

		off2 = state.OnEvent("WindowFocusTimestampChanged", func(obj interface{}) {
			ch <- ""
			off1()
			off2()
		})

		return <-ch
	}
}

func UseWorkspaceWindows(state *State) func(workspaceId int) []*Window {
	return func(workspaceId int) []*Window {
		result := make([]*Window, 0)
		windows := state.Windows

		for _, window := range windows {
			if window.WorkspaceID != workspaceId {
				continue
			}
			result = append(result, window)
		}
		slices.SortFunc(result, func(a, b *Window) int {
			ax := a.Layout.PosInScrollingLayout[0]
			ay := a.Layout.PosInScrollingLayout[1]
			bx := b.Layout.PosInScrollingLayout[0]
			by := b.Layout.PosInScrollingLayout[1]
			if ax != bx {
				return ax - bx
			}
			return ay - by
		})

		return result
	}
}

func UseUpdateSpadOriInfo(state *State) func(winId int) {
	return func(winId int) {
		state.mu.Lock()
		defer state.mu.Unlock()
		OriginWindowInfo := state.OriginWindowInfo
		Workspaces := state.Workspaces

		var workspaceId int
		for _, w := range Workspaces {
			if w.Name == nil {
				continue
			}
			if *w.Name == config.SpadWorkspaceName {
				workspaceId = w.ID
				break
			}
		}

		if workspaceId == 0 {
			return
		}

		OriginWindowInfo[winId].Workspace = workspaceId

	}
}

type ChangeWorkspaceInfo struct {
	ID       int
	Output   string
	Idx      int
	IsActive bool
}

func UseWorkspaceChangeComplete(state *State) func([]ChangeWorkspaceInfo) chan struct{} {
	return func(changes []ChangeWorkspaceInfo) chan struct{} {
		instance := GetStateInstance()
		ch := make(chan struct{}, 1)
		instance.OnEvent("WorkspacesChanged", func(i interface{}) {
			info := i.(EventStreamMsg)
			workspaces := info.WorkspacesChanged.Workspaces
			for _, item := range changes {
				matchIndex := slices.IndexFunc(workspaces, func(i Workspace) bool {
					return i.ID == item.ID
				})
				// log.Println(`test:>UseWorkspaceChangeComplete`, item.ID, matchIndex)
				if matchIndex == -1 {
					continue
				}
				matchItem := workspaces[matchIndex]
				// if matchItem.Idx != item.Idx ||
				// 	matchItem.IsActive != item.IsActive ||
				// 	matchItem.Output != item.Output {
				// 	return
				// }
				log.Println(`test:>UseWorkspaceChangeComplete`, matchItem.Output, item.Output)
				if matchItem.Output != item.Output {
					return
				}
			}

			ch <- struct{}{}
		})

		return ch
	}
}
