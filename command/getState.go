package command

import (
	"encoding/json"
	"net/http"
	"niri-startup/state"
	"niri-startup/utils"
)

func GetState(w http.ResponseWriter, r *http.Request) {
	instance := state.GetStateInstance()
	view := instance.GetSnapshot()
	str, err := json.Marshal(view)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.ReturnHttp(w, string(str))
}
