package state

import (
	"fmt"
	"sort"
	"time"
)

func UseWaitWindowOpen(state *State) func(func(*Window) bool) (*Window, error) {
	return func(filterFn func(*Window) bool) (*Window, error) {
		for _, w := range state.Windows {
			if filterFn(&w) {
				return &w, nil
			}
		}

		ch := make(chan *Window, 1)
		var off func()
		off = state.OnEvent("WindowOpenedOrChanged", func(msg interface{}) {
			w := &msg.(Msg).WindowOpenedOrChanged.Window
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
			if filterFn(&w) {
				wins = append(wins, &w)
			}
		}
		return wins
	}
}
func UseOnWindowBlur(state *State) func(*Window, func()) func() {
	return func(win *Window, fn func()) func() {
		var off func()
		off = state.OnEvent("WindowFocusChanged", func(obj interface{}) {
			if obj.(Msg).WindowFocusChanged.Id != win.ID {
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
		off1 = state.OnEvent("ScreenshotCaptured", func(obj interface{}) {
			path := obj.(Msg).ScreenshotCaptured.Path
			ch <- path
			off1()
		})

		var off2 func()
		off2 = state.OnEvent("WindowFocusTimestampChanged", func(obj interface{}) {
			ch <- ""
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
			result = append(result, &window)
		}
		sort.Slice(result, func(i, j int) bool {
			a := result[i]
			b := result[j]
			ax := a.Layout.PosInScrollingLayout[0]
			ay := a.Layout.PosInScrollingLayout[1]
			bx := b.Layout.PosInScrollingLayout[0]
			by := b.Layout.PosInScrollingLayout[1]

			return ax < bx || ay < by
		})

		return result
	}
}
