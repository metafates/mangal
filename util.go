package main

import "golang.org/x/exp/constraints"

func IfElse[T any](condition bool, a, b T) T {
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

func BytesToMegabytes(bytes int64) float64 {
	return float64(bytes) / 1_048_576 // 1024 ^ 2
}

func Min[T constraints.Ordered](a, b T) T {
	if a <= b {
		return a
	}

	return b
}

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}

	return b
}
