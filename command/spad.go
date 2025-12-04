package command

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SpadReq struct {
	Name string `json:"name"`
}

func Spad(w http.ResponseWriter, r *http.Request) {
	var req SpadReq

	// 解析 JSON
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	fmt.Println(req.Name)
}
