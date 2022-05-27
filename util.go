package main

import (
	"golang.org/x/exp/constraints"
	"os"
)

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

func PrettyTrim(text string, limit int) string {
	if len(text)-3 > limit {
		return text[:limit-3] + "..."
	}

	return text
}

// Plural makes singular word a plural if count ≠ 1
func Plural(word string, count int) string {
	if count == 1 {
		return word
	}

	return word + "s"
}

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}

	return b
}

// IsUnique tests if given array has only unique elements (e.g. if it's a set)
func IsUnique[T comparable](arr []T) bool {
	for i, a := range arr {
		for j, b := range arr {
			if i == j {
				continue
			}

			if a == b {
				return false
			}
		}
	}
	return true
}

// DirSize gets directory size in bytes
func DirSize(path string) (int64, error) {
	var size int64
	err := Afero.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}
