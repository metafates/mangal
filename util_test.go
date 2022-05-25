package main

import "testing"

func TestContains(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	validItem := 7
	invalidItem := 42

	conditions := []bool{
		Contains[int](items, validItem),
		!Contains[int](items, invalidItem),
	}

	for _, condition := range conditions {
		if !condition {
			t.Fail()
		}
	}
}

func TestIfElse(t *testing.T) {
	IfElse(true, func() {}, func() { t.Fail() })()
	IfElse(false, func() { t.Fail() }, func() {})()
}
