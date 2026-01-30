package state

import (
	"fmt"
	"niri-startup/config"
	"slices"
	"time"
)

func UseWaitWindowOpen(state *State) func(func(*Window) bool) (*Window, error) {
	return func(filterFn func(*Window) bool) (*Window, error) {
		state.mu.RLock()

		for _, w := range state.Windows {
			if filterFn(w) {
				state.mu.RUnlock()
				return w, nil
			}
		}

		state.mu.RUnlock()

		ch := make(chan *Window, 1)
		off := state.OnEvent("WindowOpenedOrChanged", func(msg interface{}) {
			w := &msg.(EventStreamMsg).WindowOpenedOrChanged.Window
			if filterFn(w) {
				ch <- w
			}
		})
		defer off()

		select {
		case w := <-ch:
			return w, nil
		case <-time.After(3 * time.Second):
			// 超时也取消事件监听
			return nil, fmt.Errorf("timeout waiting for window")
		}
	}
}

func UseWindowFilter(state *State) func(func(*Window) bool) []*Window {
	return func(filterFn func(*Window) bool) []*Window {
		state.mu.RLock()
		defer state.mu.RUnlock()
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
		state.mu.RLock()
		defer state.mu.RUnlock()
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

		owi, ok := OriginWindowInfo[winId]
		if !ok {
			return
		}
		owi.Workspace = workspaceId

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
		ch := make(chan struct{}, 1)
		done := make(chan struct{})

		offFn := state.OnEvent("localWorkspacesChanged", func(i interface{}) {
			workspaces := i.(map[int]*Workspace)

			for _, item := range changes {
				matchItem, ok := workspaces[item.ID]
				if !ok {
					continue
				}

				if matchItem.Idx != item.Idx ||
					matchItem.Output != item.Output ||
					item.IsActive != matchItem.IsActive {
					return
				}

				if item.IsActive != matchItem.IsActive {
					return
				}
			}

			close(done)
		})

		go func() {
			select {
			case <-done:
				close(ch)
			case <-time.After(2 * time.Second):
				close(ch)
			}
			offFn()
		}()

		return ch
	}
}
