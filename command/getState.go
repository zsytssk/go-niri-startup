package command

import (
	"encoding/json"
	"net/http"
	"niri-startup/state"
	"niri-startup/utils"
)

func GetState(w http.ResponseWriter, r *http.Request) {
	instance := state.GetStateInstance()
	str, err := json.Marshal(instance.CurrentWorkspaceId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.ReturnHttp(w, string(str))
}
