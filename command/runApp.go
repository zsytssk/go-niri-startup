package command

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RunAppReq struct {
	Name string `json:"name"`
}

func RunApp(w http.ResponseWriter, r *http.Request) {
	var req RunAppReq

	// 解析 JSON
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	fmt.Println(req.Name)
}
