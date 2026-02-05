package action

import (
	"encoding/json"
	"fmt"
	"niri-startup/state"
	"niri-startup/utils"
	"slices"
)

func testAction() {
	outputsStr, err := utils.NiriSend(map[string]interface{}{
		"Outputs": nil,
	})
	if err != nil {
		return
	}

	var monitors map[string]map[string]map[string]state.Monitor
	if err := json.Unmarshal([]byte(outputsStr), &monitors); err != nil {
		panic(err)
	}
	outputs := monitors["Ok"]["Outputs"]
	var arr = make([]state.Monitor, 0, len(outputs))
	for _, m := range outputs {
		arr = append(arr, m)
	}
	slices.SortFunc(arr, func(a, b state.Monitor) int {
		if a.Logical.X != b.Logical.X {
			return a.Logical.X - b.Logical.X
		}
		return a.Logical.Y - b.Logical.Y
	})

	for _, m := range arr {
		fmt.Println(m.Name)
	}

}
