package main

func IfElse[T any](condition bool, a T, b T) T {
	if condition {
		return a
	}

	return b
}

func Contains[T comparable](slice []T, elem T) bool {
	for _, el := range slice {
		if el == elem {
			return true
		}
	}

	return false
}
