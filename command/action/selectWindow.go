package action

import (
	"fmt"
	"niri-startup/state"
	"niri-startup/utils"
	"slices"
	"sort"
	"strconv"
	"strings"
)

type WindowInfo struct {
	Title     string
	ID        int
	AppId     string
	Workspace int
	OutputIdx int
	Output    string
	Idx       int
}

func SelectWindow() error {
	instance := state.GetStateInstance()
	workspace := instance.Workspaces
	windows := instance.Windows
	outputs := instance.Outputs

	wins := make([]WindowInfo, 0)
	for _, w := range windows {
		workspace, ok := workspace[w.WorkspaceID]
		if !ok {
			continue
		}
		outputIdx := slices.Index(outputs, workspace.Output)
		wins = append(wins, WindowInfo{
			Title:     w.Title,
			ID:        w.ID,
			AppId:     w.AppId,
			Workspace: workspace.Idx,
			OutputIdx: outputIdx,
			Output:    workspace.Output,
			Idx:       w.Layout.PosInScrollingLayout[0],
		})
	}
	sort.Slice(wins, func(i, j int) bool {
		a := wins[i]
		b := wins[j]
		return a.Output < b.Output || a.Workspace < b.Workspace || a.Idx < b.Idx
	})

	var lines []string
	for i, item := range wins {
		lines = append(lines,
			fmt.Sprintf("%d. %s(%s:%d:%d)",
				i+1, item.Title, item.Output, item.Workspace, item.Idx),
		)
	}
	input := strings.Join(lines, "\n")

	result, err := utils.RunCMD(fmt.Sprintf(`echo "%s" | fuzzel -d -p "请选择: "`, input), false)
	if err != nil {
		return err
	}

	parts := strings.Split(result, ".")
	if len(parts) == 0 {
		return fmt.Errorf("cant find match window")
	}

	// 1. 解析 index
	idx, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return fmt.Errorf("Error")
	}
	index := idx - 1
	// 3. 获取对应 window
	window := wins[index]

	// window 就是你要的
	utils.NiriSendActionArr([]utils.Action{
		{
			FocusWindow: &utils.WindowWithId{Id: window.ID},
		},
		{
			CenterWindow: &utils.WindowWithId{Id: window.ID},
		},
	})
	return nil
}
