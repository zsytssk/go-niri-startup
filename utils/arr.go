package utils


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



func Contains[T comparable](s []T, target T) bool {
    for _, item := range s {
        if item == target {
            return true
        }
    }
    return false
}