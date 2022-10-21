package collection

func Index[S comparable](source []S, t S) int {
	for i, v := range source {
		if v == t {
			return i
		}
	}
	return -1
}

func Include[S comparable](source []S, t S) bool {
	return Index(source, t) >= 0
}

func Any[S interface{}](source []S, f func(S) bool) bool {
	for _, v := range source {
		if f(v) {
			return true
		}
	}
	return false
}

func All[S interface{}](source []S, f func(S) bool) bool {
	for _, v := range source {
		if !f(v) {
			return false
		}
	}
	return true
}

func Filter[S interface{}](source []S, f func(S) bool) []S {
	sourcef := make([]S, 0)
	for _, v := range source {
		if f(v) {
			sourcef = append(sourcef, v)
		}
	}
	return sourcef
}

func Map[S interface{}, T interface{}](source []S, f func(S) T) []T {
	mappedArray := make([]T, len(source))
	for i, v := range source {
		mappedArray[i] = f(v)
	}
	return mappedArray
}
