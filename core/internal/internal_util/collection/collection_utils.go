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

func Count[S interface{}](source []S, f func(S) bool) int {
	count := 0

	for _, v := range source {
		if f(v) {
			count++
		}
	}

	return count
}

func Reduce[S interface{}, T interface{}](source []S, f func(S, T) T, initial T) T { //nolint:ireturn,lll // generic collection function
	accumulator := initial

	for _, v := range source {
		accumulator = f(v, accumulator)
	}

	return accumulator
}

func Prepend[S interface{}](source []S, elements ...S) []S {
	return append(elements, source...)
}

func Append[S interface{}](source []S, elements ...S) []S {
	return append(source, elements...)
}
