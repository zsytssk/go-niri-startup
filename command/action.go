package command

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ActionReq struct {
	Name string `json:"name"`
}

func Action(w http.ResponseWriter, r *http.Request) {
	var req ActionReq

	// 解析 JSON
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	fmt.Println(req.Name)
}
