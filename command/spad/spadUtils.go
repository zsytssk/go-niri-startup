package spad

import (
	"niri-startup/config"
	"niri-startup/state"
)

func UseMatchFn(spad *config.Spad) func(*state.Window) bool {
	return func(w *state.Window) bool {
		if spad.AppId != "" {
			return w.AppId == spad.AppId
		}
		return false
	}
}

func isSpadWinAction(w *state.Window) bool {
	return w.IsFloating && w.IsFocused
}
