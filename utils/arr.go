package utils

import (
	"strings"
)

func ArrJoin[T any](arr []T, fn func(item T, index int) string) string {
	var fzfInput strings.Builder

	for index, item := range arr {
		fzfInput.WriteString(fn(item, index))

	}

	return fzfInput.String()
}
