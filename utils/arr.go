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

func CloneMap[K comparable, V any](src map[K]V) map[K]V {
	if src == nil {
		return nil
	}
	dst := make(map[K]V, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
