package arrays

func Last[T any](s []T) T {
	return s[len(s)-1]
}

func Contains[T comparable](s []T, v T) bool {
	for _, item := range s {
		if item == v {
			return true
		}
	}

	return false
}
